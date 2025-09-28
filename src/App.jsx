import { BrowserRouter as Router, Routes, Route, Navigate, useLocation } from 'react-router-dom';
import { useAuthStore } from './store/authStore';
import BottomNavbar from './components/BottomNavbar';
import Home from './pages/Home';
import Checkout from './pages/Checkout';
import Transactions from './pages/Transactions';
import TransactionDetail from './pages/TransactionDetail';
import Profile from './pages/Profile';
import Login from './pages/Login';

// Admin Components
import AdminLogin from './pages/admin/AdminLogin';
import AdminLayout from './components/admin/AdminLayout';
import Dashboard from './pages/admin/Dashboard';
import AdminUsers from './pages/admin/AdminUsers';
import Categories from './pages/admin/Categories';
import Products from './pages/admin/Products';
import Vouchers from './pages/admin/Vouchers';

// 404 Components
import NotFound from './components/customer/NotFound';
import AdminNotFound from './components/admin/AdminNotFound';
import './styles/globals.css';

function AppContent() {
  const { isLoggedIn } = useAuthStore();
  const location = useLocation();
  const isAdminRoute = location.pathname.startsWith('/admin');

  if (isAdminRoute) {
    return (
      <Routes>
        <Route path="/admin/login" element={<AdminLogin />} />
        <Route path="/admin/" element={<AdminLayout><Dashboard /></AdminLayout>} />
        <Route path="/admin/users" element={<AdminLayout><AdminUsers /></AdminLayout>} />
        <Route path="/admin/categories" element={<AdminLayout><Categories /></AdminLayout>} />
        <Route path="/admin/products" element={<AdminLayout><Products /></AdminLayout>} />
        <Route path="/admin/vouchers" element={<AdminLayout><Vouchers /></AdminLayout>} />
        <Route path="/admin/*" element={<AdminNotFound />} />
      </Routes>
    );
  }

  return (
    <div className="min-h-screen bg-background text-text-primary">
      <main className="pb-16"> {/* Padding bottom untuk navbar */}
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/checkout/:productId" element={<Checkout />} />
          <Route path="/transactions" element={isLoggedIn ? <Transactions /> : <Navigate to="/login" />} />
          <Route path="/transactions/:id" element={isLoggedIn ? <TransactionDetail /> : <Navigate to="/login" />} />
          <Route path="/profile" element={isLoggedIn ? <Profile /> : <Navigate to="/login" />} />
          <Route path="/login" element={<Login />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </main>
      <BottomNavbar />
    </div>
  );
}

function App() {
  return (
    <Router>
      <AppContent />
    </Router>
  );
}

export default App;