// Based on https://github.com/inconshreveable/go-update/blob/7a872911e5b39953310f0a04161f0d50c7e63755/apply.go

package update

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// apply 使用给定的 io.Reader 的内容来更新 targetPath 的可执行文件。
//
// apply 执行以下操作以确保安全的跨平台更新：
//
// 1. 创建新文件 /path/to/target.new，并将更新文件的内容写入其中
//
// 2. 将 /path/to/target 重命名为 /path/to/target.old
//
// 3. 将 /path/to/target.new 重命名为 /path/to/target
//
// 4.如果最终的重命名成功，删除 /path/to/target.old 并返回无错误。
//
// 5. 如果最终重命名失败，尝试通过将 /path/to/target.old 重命名会
// /path/to/target 进行回滚。
//
// 如果回滚操作失败，文件系统将处于不一致状态（第 4 步和第 5 步之间），
// 既没有新的可执行文件，并且旧的可执行文件无法移动回其原始位置。在这种情况下，
// 应该通知用户这个坏消息，并要求他们手动恢复。
func apply(update io.Reader, targetPath string) error {
	newBytes, err := io.ReadAll(update)
	if err != nil {
		return err
	}

	// 获取可执行文件所在的目录
	updateDir := filepath.Dir(targetPath)
	filename := filepath.Base(targetPath)

	// 将新二进制的内容复制到新可执行文件中。
	newPath := filepath.Join(updateDir, fmt.Sprintf("%s.new", filename))
	fp, err := os.OpenFile(newPath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = io.Copy(fp, bytes.NewReader(newBytes))
	if err != nil {
		return err
	}

	// 如果我们不调用 fp.Close()，Windows 将不允许我们移动新可执行文件。
	// 因为文件仍处于 "in use"（使用中）状态。
	fp.Close()

	// 这是我们将要移动可执行文件的位置，以便可以将更新的文件替代进来
	oldPath := filepath.Join(updateDir, fmt.Sprintf("%s.old", filename))

	// 删除任何现有的就执行文件 - 这在 Windows 上是必要的，原因有两个：
	// 1. 成功更新后，Windows 无法删除 .old 文件，因为进程仍在运行
	// 2. 如果目标文件已存在，Windows 重命名操作将失败
	_ = os.Remove(oldPath)

	// 将现有的可执行文件移到同一目录下的新文件中
	err = os.Rename(targetPath, oldPath)
	if err != nil {
		return err
	}

	// 将新可执行文件移到目标位置
	err = os.Rename(newPath, targetPath)

	if err != nil {
		// 移动失败
		//
		// 文件系统现在处于不良状态。我们已成功将现有的二进制文件移动到新位置，
		// 但无法将新二进制文件移动到原来的位置。这意味着当前可执行文件的位置上没有文件！
		// 尝试通过将旧二进制文件恢复到原始路径来回滚。
		rerr := os.Rename(oldPath, targetPath)
		if rerr != nil {
			return err
		}

		return err
	}

	// 移动成功，删除旧的二进制文件
	err = os.Remove(oldPath)
	if err != nil {
		if runtime.GOOS == "windows" {
			// Windows 无法删除 .old 文件，因为进程仍在运行。删除会提示 "Access is denied"。
			// 因此，启动外部进程来删除旧的二进制文件。
			// 外部进程会等待一会以确保进程已退出。
			//
			// https://stackoverflow.com/a/73585620
			exec.Command("cmd.exe", "/c", "ping 127.0.0.1 -n 2 > NUL & del "+oldPath).Start()
			return nil
		}

		return err
	}

	return nil
}
