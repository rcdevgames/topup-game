import { Link, useLocation } from 'react-router-dom';
import { useAdminStore } from '../../store/adminStore';

function AdminSidebar() {
  const location = useLocation();
  const { sidebarOpen, admin } = useAdminStore();

  const menuItems = [
    {
      path: '/admin/',
      label: 'Dashboard',
      icon: 'fas fa-tachometer-alt',
      roles: ['super_admin', 'operator']
    },
    {
      path: '/admin/users',
      label: 'Admin Users',
      icon: 'fas fa-users-cog',
      roles: ['super_admin']
    },
    {
      path: '/admin/categories',
      label: 'Kategori Produk',
      icon: 'fas fa-tags',
      roles: ['super_admin', 'operator']
    },
    {
      path: '/admin/products',
      label: 'Produk',
      icon: 'fas fa-box',
      roles: ['super_admin', 'operator']
    },
    {
      path: '/admin/vouchers',
      label: 'Voucher',
      icon: 'fas fa-ticket-alt',
      roles: ['super_admin', 'operator']
    },
    {
      path: '/admin/transactions',
      label: 'Transaksi',
      icon: 'fas fa-credit-card',
      roles: ['super_admin', 'operator']
    }
  ];

  const filteredMenuItems = menuItems.filter(item => 
    item.roles.includes(admin.role)
  );

  return (
    <div className={`bg-primary-dark text-white transition-all duration-300 flex flex-col
      fixed left-0 top-0 h-full z-50 lg:relative lg:h-full
      ${sidebarOpen ? 'w-64' : 'w-16'}
      ${!sidebarOpen ? '-translate-x-full lg:translate-x-0' : 'translate-x-0'}
    `}>
      {/* Logo - Fixed */}
      <div className="flex-shrink-0 p-4 border-b border-primary-medium">
        {sidebarOpen ? (
          <div className="flex items-center space-x-3">
            <img 
              src="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHZpZXdCb3g9IjAgMCA0MCA0MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHJlY3QgeD0iNSIgeT0iMTAiIHdpZHRoPSIzMCIgaGVpZ2h0PSIyMCIgcng9IjMiIGZpbGw9IiM3OEEwODMiLz4KPGNpcmNsZSBjeD0iMTMiIGN5PSIxOCIgcj0iMiIgZmlsbD0id2hpdGUiLz4KPGNpcmNsZSBjeD0iMjciIGN5PSIxOCIgcj0iMiIgZmlsbD0id2hpdGUiLz4KPHJlY3QgeD0iMTUiIHk9IjMwIiB3aWR0aD0iMTAiIGhlaWdodD0iNSIgcng9IjIiIGZpbGw9IiM3OEEwODMiLz4KPC9zdmc+Cg=="
              alt="Logo"
              className="w-8 h-8"
            />
            <div>
              <h2 className="text-lg font-bold">Waw Store</h2>
              <p className="text-xs text-gray-300">Admin Panel</p>
            </div>
          </div>
        ) : (
          <div className="flex justify-center">
            <img 
              src="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDAiIGhlaWdodD0iNDAiIHZpZXdCb3g9IjAgMCA0MCA0MCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHJlY3QgeD0iNSIgeT0iMTAiIHdpZHRoPSIzMCIgaGVpZ2h0PSIyMCIgcng9IjMiIGZpbGw9IiM3OEEwODMiLz4KPGNpcmNsZSBjeD0iMTMiIGN5PSIxOCIgcj0iMiIgZmlsbD0id2hpdGUiLz4KPGNpcmNsZSBjeD0iMjciIGN5PSIxOCIgcj0iMiIgZmlsbD0id2hpdGUiLz4KPHJlY3QgeD0iMTUiIHk9IjMwIiB3aWR0aD0iMTAiIGhlaWdodD0iNSIgcng9IjIiIGZpbGw9IiM3OEEwODMiLz4KPC9zdmc+Cg=="
              alt="Logo"
              className="w-8 h-8"
            />
          </div>
        )}
      </div>

      {/* Navigation - Scrollable */}
      <nav className="flex-1 overflow-y-auto mt-6">
        {filteredMenuItems.map((item) => {
          const isActive = location.pathname === item.path;
          return (
            <Link
              key={item.path}
              to={item.path}
              className={`flex items-center px-4 py-3 text-sm transition-colors hover:bg-primary-medium ${
                isActive ? 'bg-primary-medium border-r-4 border-success' : ''
              }`}
              title={!sidebarOpen ? item.label : ''}
            >
              <i className={`${item.icon} w-5 text-center`}></i>
              {sidebarOpen && (
                <span className="ml-3">{item.label}</span>
              )}
            </Link>
          );
        })}
      </nav>

      {/* User Info - Fixed at Bottom */}
      <div className="flex-shrink-0 p-4 border-t border-primary-medium">
        {sidebarOpen ? (
          <div className="flex items-center space-x-3">
            <div className="w-8 h-8 bg-success rounded-full flex items-center justify-center">
              <i className="fas fa-user text-white text-sm"></i>
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-sm font-medium truncate">{admin.name}</p>
              <p className="text-xs text-gray-300 capitalize">{admin.role.replace('_', ' ')}</p>
            </div>
          </div>
        ) : (
          <div className="flex justify-center">
            <div className="w-8 h-8 bg-success rounded-full flex items-center justify-center">
              <i className="fas fa-user text-white text-sm"></i>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default AdminSidebar;