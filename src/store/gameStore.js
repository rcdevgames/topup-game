import { create } from 'zustand';

// Sample dummy data
const sampleProducts = [
  {
    id: 1,
    name: 'Mobile Legends 86 Diamond',
    description: '86 Diamond untuk Mobile Legends',
    price: '20000',
    image: 'https://images.unsplash.com/photo-1550745165-9bc0b252726f?w=150&h=100&fit=crop&crop=center',
    category: 'Mobile Legends'
  },
  {
    id: 2,
    name: 'Free Fire 70 Diamond',
    description: '70 Diamond untuk Free Fire',
    price: '15000',
    image: 'https://images.unsplash.com/photo-1538481199705-c710c4e965fc?w=150&h=100&fit=crop&crop=center',
    category: 'Free Fire'
  },
  {
    id: 3,
    name: 'Mobile Legends 172 Diamond',
    description: '172 Diamond untuk Mobile Legends',
    price: '40000',
    image: 'https://images.unsplash.com/photo-1556438064-2d7646166914?w=150&h=100&fit=crop&crop=center',
    category: 'Mobile Legends'
  },
  {
    id: 4,
    name: 'Free Fire 140 Diamond',
    description: '140 Diamond untuk Free Fire',
    price: '30000',
    image: 'https://images.unsplash.com/photo-1542751371-adc38448a05e?w=150&h=100&fit=crop&crop=center',
    category: 'Free Fire'
  },
];

const sampleTransactions = [
  {
    id: 'TRX001',
    game: 'Mobile Legends',
    amount: '20000',
    status: 'success',
    date: '2024-03-15'
  },
  {
    id: 'TRX002',
    game: 'Free Fire',
    amount: '15000',
    status: 'pending',
    date: '2024-03-16'
  }
];

export const useGameStore = create((set, get) => ({
  products: sampleProducts,
  allProducts: sampleProducts, // Keep original for search
  categories: [
    { id: 1, name: 'Mobile Legends' },
    { id: 2, name: 'Free Fire' },
    { id: 3, name: 'PUBG Mobile' },
    { id: 4, name: 'Genshin Impact' }
  ],
  transactions: sampleTransactions,
  gameAccounts: [],
  cart: [],
  
  loadMoreProducts: () => {
    // Simulate loading more products
    const gameImages = [
      'https://images.unsplash.com/photo-1511512578047-dfb367046420?w=150&h=100&fit=crop&crop=center',
      'https://images.unsplash.com/photo-1493711662062-fa541adb3fc8?w=150&h=100&fit=crop&crop=center',
      'https://images.unsplash.com/photo-1552820728-8b83bb6b773f?w=150&h=100&fit=crop&crop=center',
      'https://images.unsplash.com/photo-1509198397868-475647b2a1e5?w=150&h=100&fit=crop&crop=center'
    ];
    const categories = ['Mobile Legends', 'Free Fire', 'PUBG Mobile', 'Genshin Impact'];
    const randomCategory = categories[Math.floor(Math.random() * categories.length)];
    const randomImage = gameImages[Math.floor(Math.random() * gameImages.length)];
    
    const newProducts = [
      {
        id: Date.now(),
        name: `${randomCategory} Diamond Pack`,
        description: `Diamond premium untuk ${randomCategory}`,
        price: (Math.floor(Math.random() * 50) + 10) * 1000 + '',
        image: randomImage,
        category: randomCategory
      }
    ];
    set((state) => ({ products: [...state.products, ...newProducts] }));
  },
  
  searchProducts: (term) => {
    const { allProducts } = get();
    if (!term) {
      set({ products: allProducts });
      return;
    }
    const filtered = allProducts.filter(product => 
      product.name.toLowerCase().includes(term.toLowerCase()) ||
      product.category.toLowerCase().includes(term.toLowerCase())
    );
    set({ products: filtered });
  },
  
  loadMoreTransactions: () => {
    const newTransaction = {
      id: `TRX${Date.now()}`,
      game: 'Mobile Legends',
      amount: '35000',
      status: 'success',
      date: new Date().toISOString().split('T')[0]
    };
    set((state) => ({ transactions: [...state.transactions, newTransaction] }));
  },
  
  addToCart: (product) => {
    set((state) => ({ cart: [...state.cart, product] }));
  },
  
  addGameAccount: (account) => {
    const newAccount = { ...account, id: Date.now().toString() };
    set((state) => ({ gameAccounts: [...state.gameAccounts, newAccount] }));
  },
  
  updateGameAccount: (id, updated) => {
    set((state) => ({
      gameAccounts: state.gameAccounts.map(acc => acc.id === id ? updated : acc)
    }));
  },
  
  deleteGameAccount: (id) => {
    set((state) => ({
      gameAccounts: state.gameAccounts.filter(acc => acc.id !== id)
    }));
  },
  
  addTransaction: (transaction) => {
    set((state) => ({ transactions: [transaction, ...state.transactions] }));
  },
}));