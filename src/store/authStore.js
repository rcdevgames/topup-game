import { create } from 'zustand';

export const useAuthStore = create((set) => ({
  isLoggedIn: false,
  user: { name: '', phone: '', password: '' },
  login: (userData) => set({ isLoggedIn: true, user: userData }),
  logout: () => set({ isLoggedIn: false, user: {} }),
  updateUser: (newUser) => set((state) => ({ user: { ...state.user, ...newUser } })),
}));