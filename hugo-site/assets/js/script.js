document.addEventListener('DOMContentLoaded', function() {
  // Mobile menu toggle
  const menuBtn = document.getElementById('mobile-menu-btn');
  const mobileMenu = document.getElementById('mobile-menu');
  if (menuBtn && mobileMenu) {
    menuBtn.addEventListener('click', function() {
      mobileMenu.classList.toggle('hidden');
    });
  }

  // Smooth scroll for anchor links
  document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
      e.preventDefault();
      const target = document.querySelector(this.getAttribute('href'));
      if (target) target.scrollIntoView({ behavior: 'smooth' });
    });
  });

  // Hero image slideshow (background only)
  const bgSlides = document.querySelectorAll('#hero-slides .hero-slide');
  if (bgSlides.length > 1) {
    let current = 0;
    function goToSlide(n) {
      bgSlides[current].classList.replace('opacity-100', 'opacity-0');
      current = (n + bgSlides.length) % bgSlides.length;
      bgSlides[current].classList.replace('opacity-0', 'opacity-100');
    }
    setInterval(() => goToSlide(current + 1), 15000);
  }
});
