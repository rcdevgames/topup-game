import { useAdminStore } from '../../store/adminStore';
import Swal from 'sweetalert2';
import { useNavigate } from 'react-router-dom';

function AdminHeader() {
  const { toggleSidebar, sidebarOpen, admin, adminLogout } = useAdminStore();
  const navigate = useNavigate();

  const handleLogout = () => {
    Swal.fire({
      title: 'Keluar dari Admin Panel?',
      icon: 'question',
      showCancelButton: true,
      confirmButtonColor: '#EF4444',
      cancelButtonColor: '#6B7280',
      confirmButtonText: 'Ya, Keluar',
      cancelButtonText: 'Batal'
    }).then((result) => {
      if (result.isConfirmed) {
        adminLogout();
        navigate('/admin/login');
        Swal.fire({
          icon: 'success',
          title: 'Berhasil Keluar',
          timer: 1500,
          showConfirmButton: false,
        });
      }
    });
  };

  return (
    <header className="flex-shrink-0 bg-white shadow-sm border-b border-gray-200 px-4 lg:px-6 py-4">
      <div className="flex items-center justify-between">
        {/* Left Side */}
        <div className="flex items-center space-x-2 md:space-x-4">
          <button
            onClick={toggleSidebar}
            className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
            aria-label="Toggle Sidebar"
          >
            <i className={`fas ${sidebarOpen ? 'fa-times' : 'fa-bars'} text-gray-600`}></i>
          </button>
          
          <h1 className="text-lg md:text-xl font-semibold text-gray-800">
            Admin Panel
          </h1>
        </div>

        {/* Right Side */}
        <div className="flex items-center space-x-2 md:space-x-4">
          {/* Notifications */}
          <button className="p-2 hover:bg-gray-100 rounded-lg transition-colors relative">
            <i className="fas fa-bell text-gray-600"></i>
            <span className="absolute -top-1 -right-1 w-3 h-3 bg-red-500 rounded-full"></span>
          </button>

          {/* Admin Profile Dropdown */}
          <div className="relative group">
            <button className="flex items-center space-x-2 md:space-x-3 p-2 hover:bg-gray-100 rounded-lg transition-colors">
              <div className="w-8 h-8 bg-primary-light rounded-full flex items-center justify-center">
                <i className="fas fa-user text-white text-sm"></i>
              </div>
              <div className="text-left hidden md:block">
                <p className="text-sm font-medium text-gray-800">{admin.name}</p>
                <p className="text-xs text-gray-500 capitalize">{admin.role.replace('_', ' ')}</p>
              </div>
              <i className="fas fa-chevron-down text-gray-400 text-xs hidden md:block"></i>
            </button>

            {/* Dropdown Menu */}
            <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 z-50">
              <div className="py-2">
                <button
                  onClick={() => navigate('/')}
                  className="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 flex items-center space-x-2"
                >
                  <i className="fas fa-home text-gray-400"></i>
                  <span>Ke Beranda</span>
                </button>
                <button
                  onClick={handleLogout}
                  className="w-full text-left px-4 py-2 text-sm text-red-700 hover:bg-red-50 flex items-center space-x-2"
                >
                  <i className="fas fa-sign-out-alt text-red-400"></i>
                  <span>Keluar</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
}

export default AdminHeader;