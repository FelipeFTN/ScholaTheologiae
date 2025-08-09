// Simple Search Modal
document.addEventListener('DOMContentLoaded', function() {
  
  const searchIcon = document.getElementById('search-icon');
  const searchModal = document.getElementById('search-modal');
  const closeSearch = document.getElementById('close-search');
  const searchInput = document.getElementById('search-input');
  
  if (!searchIcon) {
    console.error('Search icon not found!');
    return;
  }
  
  if (!searchModal) {
    console.error('Search modal not found!');
    return;
  }
  
  // Open modal when clicking search icon
  searchIcon.addEventListener('click', function() {
    searchModal.classList.add('show');
    if (searchInput) {
      searchInput.focus();
    }
  });
  
  // Close modal when clicking close button
  if (closeSearch) {
    closeSearch.addEventListener('click', function() {
      searchModal.classList.remove('show');
      if (searchInput) {
        searchInput.value = '';
      }
    });
  }
  
  // Close modal when clicking outside
  searchModal.addEventListener('click', function(e) {
    if (e.target === searchModal) {
      searchModal.classList.remove('show');
      if (searchInput) {
        searchInput.value = '';
      }
    }
  });
  
  // Close modal with Escape key
  document.addEventListener('keydown', function(e) {
    if (e.key === 'Escape' && searchModal.classList.contains('show')) {
      searchModal.classList.remove('show');
      if (searchInput) {
        searchInput.value = '';
      }
    }
  });
  
  // Handle search input
  if (searchInput) {
    searchInput.addEventListener('input', function() {
      const query = this.value.trim();
      if (query.length > 2) {
        // TODO: Implement search functionality with your API
        console.log('Searching for:', query);
        // This is where you'll call your API: http://localhost:8080/v1/search?q="query"
      }
    });
  }
});
