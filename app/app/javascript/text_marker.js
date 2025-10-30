class TextMarker {
  constructor() {
    // Prevent multiple instances
    if (document.getElementById('text-marker-btn')) {
      return;
    }
    
    this.isMarkingMode = false;
    this.markedTexts = JSON.parse(localStorage.getItem('scholaTheologiae_markedTexts') || '{}');
    this.init();
  }

  init() {
    // Check if content area exists
    const content = document.querySelector('.book-question-content');
    if (!content) return;

    this.createButton();
    this.setupEvents();
    this.restoreMarks();
  }

  createButton() {
    const button = document.createElement('div');
    button.id = 'text-marker-btn';
    button.className = 'text-marker-button';
    button.innerHTML = `
      <svg width="30px" height="30px" viewBox="0 0 512 512"  xml:space="preserve">
      <style type="text/css"> <![CDATA[ .st0{fill:#ffffff;} ]]> </style>
      <g>
        <path class="st0" d="M208.125,399.438c0,0,9.656-15.688,21.266-34.563L75.469,270.156C63.859,289,54.188,304.719,54.188,304.719
          s27.797,60.406-9.906,121.656l29.844,18.375l29.844,18.375C141.656,401.844,208.125,399.438,208.125,399.438z"/>
        <path class="st0" d="M389.531,104.688c6.031-9.844,2.969-22.719-6.859-28.781L264.359,3.109C260.938,1,257.156,0,253.422,0
          c-7.016,0-13.875,3.531-17.828,9.953L83.188,257.594l153.938,94.719L389.531,104.688z M128.344,246.844L257.313,37.313
          l98.016,60.313L226.375,307.156L128.344,246.844z"/>
        <polygon class="st0" points="22.531,488.641 74.188,488.641 87.484,467.031 48.219,442.875 	"/>
        <path class="st0" d="M482.406,484.453H117.844c-3.906,0-7.063,3.156-7.063,7.047v13.438c0,3.906,3.156,7.063,7.063,7.063h364.563
          c3.906,0,7.063-3.156,7.063-7.063V491.5C489.469,487.609,486.313,484.453,482.406,484.453z"/>
      </g>
      </svg>
    `;
    button.addEventListener('click', () => this.toggleMode());
    document.body.appendChild(button);
  }

  setupEvents() {
    document.addEventListener('mouseup', (e) => {
      if (this.isMarkingMode) {
        setTimeout(() => this.handleSelection(), 50);
      }
    });

    // Handle clicks on marked text
    document.addEventListener('click', (e) => {
      if (e.target.classList.contains('marked-text')) {
        this.showConfirm('Remover esta marcação?', () => this.removeMark(e.target));
      }
    });
  }

  toggleMode() {
    this.isMarkingMode = !this.isMarkingMode;
    const button = document.getElementById('text-marker-btn');
    const content = document.querySelector('.book-question-content');
    
    if (this.isMarkingMode) {
      button.classList.add('active');
      content.style.userSelect = 'text';
      this.showToast('Modo de marcação ativado');
    } else {
      button.classList.remove('active');
      content.style.userSelect = 'none';
      this.showToast('Modo de marcação desativado');
    }
  }

  handleSelection() {
    const selection = window.getSelection();
    const text = selection.toString().trim();
    console.log(text)
    
    if (text.length < 3) return;
    
    const range = selection.getRangeAt(0);
    // if (range.startContainer.closest('.marked-text')) {
    //   this.showToast('Texto já marcado');
    //   return;
    // }

    const mark = document.createElement('mark');
    mark.className = 'marked-text';
    mark.setAttribute('data-id', Date.now());
    
    try {
      range.surroundContents(mark);
      this.saveToStorage(text);
      selection.removeAllRanges();
      this.showToast('Texto marcado!');
    } catch (error) {
      this.showToast('Erro ao marcar texto');
    }
  }

  removeMark(element) {
    const parent = element.parentNode;
    parent.insertBefore(document.createTextNode(element.textContent), element);
    parent.removeChild(element);
    parent.normalize();
    this.showToast('Marcação removida');
  }

  saveToStorage(text) {
    const page = window.location.pathname;
    if (!this.markedTexts[page]) this.markedTexts[page] = [];
    this.markedTexts[page].push({ text, timestamp: Date.now() });
    localStorage.setItem('scholaTheologiae_markedTexts', JSON.stringify(this.markedTexts));
  }

  restoreMarks() {
    const page = window.location.pathname;
    const marks = this.markedTexts[page] || [];
    
    marks.forEach(markData => {
      const walker = document.createTreeWalker(
        document.querySelector('.book-question-content'),
        NodeFilter.SHOW_TEXT
      );
      
      let node;
      while (node = walker.nextNode()) {
        const index = node.textContent.indexOf(markData.text);
        if (index >= 0) {
          const range = document.createRange();
          range.setStart(node, index);
          range.setEnd(node, index + markData.text.length);
          
          const mark = document.createElement('mark');
          mark.className = 'marked-text';
          mark.setAttribute('data-id', markData.timestamp);
          
          try {
            range.surroundContents(mark);
            break;
          } catch (e) {}
        }
      }
    });
  }

  showToast(message) {
    const toast = document.createElement('div');
    toast.className = 'marker-toast show';
    toast.textContent = message;
    document.body.appendChild(toast);
    
    setTimeout(() => {
      toast.classList.remove('show');
      setTimeout(() => toast.remove(), 300);
    }, 2000);
  }

  showConfirm(message, onConfirm) {
    const modal = document.createElement('div');
    modal.className = 'marker-modal';
    modal.style.opacity = '1';
    modal.innerHTML = `
      <div class="modal-content">
        <div class="modal-header"><h3>Confirmação</h3></div>
        <div class="modal-body"><p>${message}</p></div>
        <div class="modal-footer">
          <button class="modal-btn modal-btn-cancel">Cancelar</button>
          <button class="modal-btn modal-btn-confirm">Confirmar</button>
        </div>
      </div>
    `;
    
    modal.querySelector('.modal-btn-cancel').onclick = () => modal.remove();
    modal.querySelector('.modal-btn-confirm').onclick = () => {
      onConfirm();
      modal.remove();
    };
    modal.onclick = (e) => e.target === modal && modal.remove();
    
    document.body.appendChild(modal);
  }
}

// Initialize when the trigger element is found
document.addEventListener('DOMContentLoaded', () => {
  const trigger = document.querySelector('[data-text-marker="true"]');
  if (trigger && !document.getElementById('text-marker-btn')) {
    new TextMarker();
  }
});
