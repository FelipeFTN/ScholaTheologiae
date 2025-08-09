// Search modal functionality
document.addEventListener('DOMContentLoaded', function() {
  
  // Get all the elements we need
  const searchIcon = document.getElementById('search-icon');
  const searchModal = document.getElementById('search-modal');
  const closeSearch = document.getElementById('close-search');
  const searchInput = document.getElementById('search-input');
  const searchForm = document.getElementById('search-form');
  const searchButton = document.querySelector('.search-submit-btn');
  
  // Bail out if modal doesn't exist
  if (!searchModal || !searchIcon) return;
  
  // Open modal and focus input
  searchIcon.addEventListener('click', function() {
    searchModal.classList.add('show');
    searchInput?.focus();
  });
  
  // Close modal function
  function closeModal() {
    searchModal.classList.remove('show');
    if (searchInput) searchInput.value = '';
  }
  
  // Close with X button
  closeSearch?.addEventListener('click', closeModal);
  
  // Close when clicking outside
  searchModal.addEventListener('click', function(e) {
    if (e.target === searchModal) closeModal();
  });
  
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
});
