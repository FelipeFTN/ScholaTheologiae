// Simple search results page handler
function handleSearchResults() {
  // Only run on search results page
  const searchPage = document.querySelector('.search-results-page');
  if (!searchPage) return;

  // Get search query from URL  
  const searchResults = document.querySelector('.search-results-section');
  if (searchResults) {
    searchResults.style.display = 'none'; // Hide results initially
  }

  // Hide loading spinner after a quick delay for nice UX
  const loadingStatus = document.querySelector('.search-status');
  if (loadingStatus) {
    setTimeout(() => {
      loadingStatus.style.display = 'none';
      if (searchResults) {
        searchResults.style.display = 'block';
      }
    }, 500);
  }
}

// Run when page loads (works with Turbo navigation too)
document.addEventListener('DOMContentLoaded', handleSearchResults);
document.addEventListener('turbo:load', handleSearchResults);
