// Search modal functionality
function initializeSearchModal() {
  // Get all the elements we need
  const searchIcon = document.getElementById('search-icon');
  const searchModal = document.getElementById('search-modal');
  const closeSearch = document.getElementById('close-search');
  const searchInput = document.getElementById('search-input');
  const searchForm = document.getElementById('search-form');
  const searchButton = document.querySelector('.search-submit-btn');
  
  // Bail out if modal doesn't exist
  if (!searchModal || !searchIcon) return;

  // Turbo swaps the body on navigation, so a fresh modal element is wired on every
  // visit; skip the wiring when this element was already set up.
  if (searchModal.dataset.ready === 'true') return;
  searchModal.dataset.ready = 'true';
  
  // Close modal function
  function closeModal() {
    searchModal.classList.remove('show');
    if (searchInput) searchInput.value = '';
  }

  // Open modal and focus input
  searchIcon.addEventListener('click', function(e) {
    e.preventDefault();
    e.stopPropagation();
    searchModal.classList.add('show');
    searchInput?.focus();
  });

  // Close with X button - prevent event bubbling
  if (closeSearch) {
    closeSearch.addEventListener('click', function(e) {
      e.preventDefault();
      e.stopPropagation();
      closeModal();
    });
  }
  
  // Close when clicking outside - but not when clicking inside modal content
  searchModal.addEventListener('click', function(e) {
    // Only close if clicking the overlay itself, not the content
    if (e.target === searchModal || e.target.classList.contains('search-modal-overlay')) {
      closeModal();
    }
  });
  
  // Prevent modal from closing when clicking inside content
  const modalContent = searchModal.querySelector('.search-modal-content');
  if (modalContent) {
    modalContent.addEventListener('click', function(e) {
      e.stopPropagation();
    });
  }
  
  // Close with Escape key
  document.addEventListener('keydown', function(e) {
    if (e.key === 'Escape' && searchModal.classList.contains('show')) {
      closeModal();
    }
  });
  
  // Handle search form
  if (searchForm && searchInput && searchButton) {
    
    // Enable/disable button based on input length
    searchInput.addEventListener('input', function() {
      const hasEnoughText = this.value.trim().length >= 2;
      searchButton.disabled = !hasEnoughText;
      searchButton.style.opacity = hasEnoughText ? '1' : '0.6';
    });
    
    // Initially disable button
    searchButton.disabled = true;
    searchButton.style.opacity = '0.6';
    
    // Handle form submission
    searchForm.addEventListener('submit', function(e) {
      const query = searchInput.value.trim();
      
      // Block short queries
      if (query.length < 2) {
        e.preventDefault();
        searchInput.focus();
        return;
      }
      
      // Show loading state
      searchButton.innerHTML = '<i class="nf nf-fa-spinner"></i><span>Buscando...</span>';
      searchButton.disabled = true;
    });
  }
}

// Initialize on page load and on Turbo navigation
document.addEventListener('DOMContentLoaded', initializeSearchModal);
document.addEventListener('turbo:load', initializeSearchModal);
