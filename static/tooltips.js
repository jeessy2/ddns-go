class Tooltip {
  constructor(element, triggers) {
    this.$element = element;
    this.$tooltip = null;
    this.originalTitle = '';
    this._bindEvents(triggers);
  }

  _createTooltipElement(options) {
    const title = options.title || this.$element.dataset.title || this.originalTitle;
    if (!title) {
      return;
    }
    const useHtml = options.hasOwnProperty('html') ? options.html : this.$element.dataset.html === 'true';
    let placement = options.placement || this.$element.dataset.placement || 'auto';
    if (placement === 'auto') {
      const rect = this.$element.getBoundingClientRect();
      const viewportWidth = window.innerWidth || document.documentElement.clientWidth;
      const viewportHeight = window.innerHeight || document.documentElement.clientHeight;
      const space = {
        top: rect.top,
        bottom: viewportHeight - rect.bottom,
        left: rect.left,
        right: viewportWidth - rect.right
      };
      placement = Object.keys(space).reduce((a, b) => space[a] > space[b] ? a : b);
    }
    this.$tooltip = html2Element(`
      <div class="tooltip bs-tooltip-${placement}"
        x-placement="${placement}"
        style="will-change: transform;"
        role="tooltip"
      >
        <div class="arrow"></div>
        <div class="tooltip-inner"></div>
      </div>
    `)
    if (useHtml) {
      this.$tooltip.querySelector('.tooltip-inner').innerHTML = title
    } else {
      this.$tooltip.querySelector('.tooltip-inner').textContent = title
    }
  }

  _updatePosition() {
    const elRect = this.$element.getBoundingClientRect()
    const bodyRect = document.body.getBoundingClientRect()
    const tooltipRect = this.$tooltip.getBoundingClientRect()
    const placement = this.$tooltip.getAttribute('x-placement')

    let left, top;
    
    switch(placement) {
      case 'top':
        left = elRect.left + (elRect.width - tooltipRect.width) / 2
        top = elRect.top - tooltipRect.height - 8
        break
      case 'bottom':
        left = elRect.left + (elRect.width - tooltipRect.width) / 2
        top = elRect.bottom + 8
        break
      case 'left':
        left = elRect.left - tooltipRect.width - 8
        top = elRect.top + (elRect.height - tooltipRect.height) / 2
        break
      case 'right':
        left = elRect.right + 8
        top = elRect.top + (elRect.height - tooltipRect.height) / 2
        break
    }

    // 考虑滚动条的影响
    left = left - bodyRect.left
    top = top - bodyRect.top
    
    this.$tooltip.style.left = `${left}px`
    this.$tooltip.style.top = `${top}px`
  }

  async show(options = {}) {
    if (this.$tooltip) {
      this.$tooltip.remove();
    }
    if (this.$element.title) {
      this.originalTitle = this.$element.title;
      this.$element.title = '';
    }
    this._createTooltipElement(options);
    if (!this.$tooltip) {
      return;
    }
    document.body.appendChild(this.$tooltip);
    await delay(0);
    if (!this.$tooltip) {
      return;
    }
    this._updatePosition();
    this.$tooltip.classList.add('show');
  }

  async hide() {
    if (this.originalTitle && !this.$element.title) {
      this.$element.title = this.originalTitle;
    }
    if (!this.$tooltip) {
      return;
    }
    this.$tooltip.classList.remove('show');
    await delay(200);
    if (!this.$tooltip) {
      return;
    }
    this.$tooltip.remove();
    this.$tooltip = null;
  }

  _bindEvents(triggers) {
    let state = 0;
    const _enter = () => {
      state += 1;
      this.show();
    };
    const _leave = () => {
      state -= 1;
      if (state <= 0) {
        this.hide();
      }
    };
    if (!triggers) {
      triggers = (this.$element.dataset.trigger || 'hover focus').split(' ');
    }
    triggers.forEach(trigger => {
      switch(trigger) {
        case 'hover':
          this.$element.addEventListener('mouseenter', _enter);
          this.$element.addEventListener('mouseleave', _leave);
          break;
        case 'focus':
          this.$element.addEventListener('focusin', _enter);
          this.$element.addEventListener('focusout', _leave);
          break;
        case 'click':
          this.$element.addEventListener('click', () => {
            if (this.$tooltip) {
              this.hide();
            } else {
              this.show();
            }
          });
          break;
        case 'manual':
          break;
        default:
          console.warn(`Unknown trigger: ${trigger}`);
      }
    });
  }
}

// 初始化所有带data-tooltip属性的元素
const initTooltips = () => {
  window.tooltips = {};
  document.querySelectorAll('[data-toggle="tooltip"]').forEach(element => {
    let key = element.dataset.tooltipKey || element.id;
    if (!key) {
      key = crypto.randomUUID();
      element.dataset.tooltipKey = key;
    }
    window.tooltips[key] = new Tooltip(element);
  });
};

// 页面加载完成后初始化
document.addEventListener('DOMContentLoaded', initTooltips); 