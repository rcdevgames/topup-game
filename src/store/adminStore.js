import { create } from 'zustand';

export const useAdminStore = create((set, get) => ({
  // Admin Auth
  isAdminLoggedIn: false,
  admin: { id: '', username: '', name: '', role: '' },
  
  // Admin Login/Logout
  adminLogin: (adminData) => set({ isAdminLoggedIn: true, admin: adminData }),
  adminLogout: () => set({ isAdminLoggedIn: false, admin: {} }),
  
  // Sidebar State
  sidebarOpen: true,
  toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
  setSidebarOpen: (open) => set({ sidebarOpen: open }),
  
  // Dashboard Analytics
  analytics: {
    totalSales: 0,
    totalTransactions: 0,
    totalUsers: 0,
    totalProducts: 0,
    dailySales: [],
    topProducts: [],
    recentTransactions: []
  },
  updateAnalytics: (data) => set((state) => ({ 
    analytics: { ...state.analytics, ...data } 
  })),
  
  // Admin Users CRUD
  adminUsers: [
    { id: 1, username: 'admin', name: 'Super Admin', role: 'super_admin', status: 'active', createdAt: '2025-09-26' },
    { id: 2, username: 'operator', name: 'Operator 1', role: 'operator', status: 'active', createdAt: '2025-09-26' }
  ],
  addAdminUser: (user) => set((state) => ({ 
    adminUsers: [...state.adminUsers, { ...user, id: Date.now() }] 
  })),
  updateAdminUser: (id, updates) => set((state) => ({
    adminUsers: state.adminUsers.map(user => user.id === id ? { ...user, ...updates } : user)
  })),
  deleteAdminUser: (id) => set((state) => ({
    adminUsers: state.adminUsers.filter(user => user.id !== id)
  })),
  
  // Categories CRUD
  categories: [
    { id: 1, name: 'Mobile Legends', description: 'MOBA Game', status: 'active', createdAt: '2025-09-26' },
    { id: 2, name: 'Free Fire', description: 'Battle Royale', status: 'active', createdAt: '2025-09-26' },
    { id: 3, name: 'PUBG Mobile', description: 'Battle Royale', status: 'active', createdAt: '2025-09-26' },
    { id: 4, name: 'Genshin Impact', description: 'RPG Adventure', status: 'active', createdAt: '2025-09-26' }
  ],
  addCategory: (category) => set((state) => ({ 
    categories: [...state.categories, { ...category, id: Date.now(), createdAt: new Date().toISOString().split('T')[0] }] 
  })),
  updateCategory: (id, updates) => set((state) => ({
    categories: state.categories.map(cat => cat.id === id ? { ...cat, ...updates } : cat)
  })),
  deleteCategory: (id) => set((state) => ({
    categories: state.categories.filter(cat => cat.id !== id)
  })),
  
  // Products CRUD with Dynamic Form Config
  products: [
    {
      id: 1,
      name: '50 Diamonds',
      categoryId: 1,
      category: 'Mobile Legends',
      price: '15000',
      description: '50 Diamonds Mobile Legends',
      image: 'https://images.unsplash.com/photo-1511512578047-dfb367046420?w=400&h=300&fit=crop&crop=center',
      status: 'active',
      formConfig: [
        { field: 'gameAccount', label: 'User ID', type: 'text', required: true, placeholder: 'Masukkan User ID' },
        { field: 'gameZone', label: 'Zone ID', type: 'text', required: true, placeholder: 'Masukkan Zone ID' }
      ],
      createdAt: '2025-09-26'
    }
  ],
  addProduct: (product) => set((state) => ({ 
    products: [...state.products, { ...product, id: Date.now(), createdAt: new Date().toISOString().split('T')[0] }] 
  })),
  updateProduct: (id, updates) => set((state) => ({
    products: state.products.map(product => product.id === id ? { ...product, ...updates } : product)
  })),
  deleteProduct: (id) => set((state) => ({
    products: state.products.filter(product => product.id !== id)
  })),
  
  // Vouchers CRUD
  vouchers: [
    {
      id: 1,
      code: 'NEWUSER10',
      type: 'percentage',
      value: 10,
      description: 'Diskon 10% untuk user baru',
      applicationType: 'all', // 'all', 'category', 'product'
      applicableIds: [], // category IDs or product IDs
      quota: 100,
      usedCount: 25,
      startDate: '2025-09-01',
      endDate: '2025-12-31',
      status: 'active',
      createdAt: '2025-09-26'
    },
    {
      id: 2,
      code: 'SAVE5K',
      type: 'fixed',
      value: 5000,
      description: 'Potongan Rp 5.000',
      applicationType: 'category',
      applicableIds: [1], // Mobile Legends only
      quota: 50,
      usedCount: 12,
      startDate: '2025-09-15',
      endDate: '2025-10-15',
      status: 'active',
      createdAt: '2025-09-26'
    }
  ],
  addVoucher: (voucher) => set((state) => ({ 
    vouchers: [...state.vouchers, { ...voucher, id: Date.now(), usedCount: 0, createdAt: new Date().toISOString().split('T')[0] }] 
  })),
  updateVoucher: (id, updates) => set((state) => ({
    vouchers: state.vouchers.map(voucher => voucher.id === id ? { ...voucher, ...updates } : voucher)
  })),
  deleteVoucher: (id) => set((state) => ({
    vouchers: state.vouchers.filter(voucher => voucher.id !== id)
  })),
  
  // Load Analytics Data
  loadDashboardData: () => {
    // Simulate API call
    const { products, vouchers } = get();
    const totalProducts = products.length;
    const totalVouchers = vouchers.length;
    
    // Mock data for analytics dengan angka besar untuk test formatting
    const analytics = {
      totalSales: 2450000000, // 2.45 Miliar
      totalTransactions: 125000, // 125 ribu
      totalUsers: 15000, // 15 ribu
      totalProducts,
      dailySales: [
        { date: '2025-09-20', sales: 250000000 },
        { date: '2025-09-21', sales: 320000000 },
        { date: '2025-09-22', sales: 280000000 },
        { date: '2025-09-23', sales: 410000000 },
        { date: '2025-09-24', sales: 350000000 },
        { date: '2025-09-25', sales: 420000000 },
        { date: '2025-09-26', sales: 470000000 }
      ],
      topProducts: [
        { name: '50 Diamonds ML', sales: 1500, revenue: 22500000 },
        { name: '100 Diamonds ML', sales: 1200, revenue: 36000000 },
        { name: '500 UC PUBG', sales: 800, revenue: 40000000 }
      ],
      recentTransactions: [
        { id: 'TRX001', user: 'User123', product: '50 Diamonds ML', amount: 15000, status: 'success' },
        { id: 'TRX002', user: 'User456', product: '100 UC FF', amount: 12000, status: 'pending' },
        { id: 'TRX003', user: 'User789', product: '300 Diamonds ML', amount: 85000, status: 'success' }
      ]
    };
    
    set({ analytics });
  }
}));