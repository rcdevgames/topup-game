import { Link, useLocation } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';

function BottomNavbar() {
  const location = useLocation();
  const { isLoggedIn } = useAuthStore();

  const menus = [
    { path: '/', label: 'Beranda', icon: 'fas fa-home' },
    { path: '/transactions', label: 'Transaksi', icon: 'fas fa-credit-card', requiresAuth: true },
    { path: '/profile', label: 'Profil', icon: 'fas fa-user', requiresAuth: true },
  ];

  return (
    <nav className="fixed bottom-0 left-0 right-0 bg-primary-dark shadow-lg z-50">
      <div className="flex justify-around py-2">
        {menus.map(({ path, label, icon, requiresAuth }) => {
          if (requiresAuth && !isLoggedIn) {
            return (
              <Link
                key={path}
                to="/login"
                className="flex flex-col items-center p-2 rounded-lg transition-colors text-text-secondary hover:text-white focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                aria-label={label}
              >
                <i className={`${icon} text-xl`}></i>
                <span className="text-xs mt-1">{label}</span>
              </Link>
            );
          }
          
          const isActive = location.pathname === path;
          return (
            <Link
              key={path}
              to={path}
              className={`flex flex-col items-center p-2 rounded-lg transition-colors ${
                isActive ? 'bg-primary-medium text-white' : 'text-text-secondary hover:text-white'
              } focus:ring-2 focus:ring-offset-2 focus:ring-blue-500`}
              aria-label={label}
            >
              <i className={`${icon} text-xl`}></i>
              <span className="text-xs mt-1">{label}</span>
            </Link>
          );
        })}
      </div>
    </nav>
  );
}

export default BottomNavbar;