import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { useNavigate } from 'react-router-dom';
import { useAdminStore } from '../../store/adminStore';
import Swal from 'sweetalert2';

const adminLoginSchema = yup.object({
  username: yup.string().required('Username wajib diisi'),
  password: yup.string().required('Password wajib diisi'),
});

function AdminLogin() {
  const navigate = useNavigate();
  const { adminLogin } = useAdminStore();
  
  const { register, handleSubmit, formState: { errors } } = useForm({
    resolver: yupResolver(adminLoginSchema),
  });

  const onSubmit = (data) => {
    // Demo admin login - in real app, this would call an API
    if ((data.username === 'admin' && data.password === 'admin123') || 
        (data.username === 'operator' && data.password === 'operator123')) {
      
      const adminData = {
        id: data.username === 'admin' ? 1 : 2,
        username: data.username,
        name: data.username === 'admin' ? 'Super Admin' : 'Operator 1',
        role: data.username === 'admin' ? 'super_admin' : 'operator'
      };
      
      adminLogin(adminData);
      Swal.fire({
        icon: 'success',
        title: 'Login Admin Berhasil',
        text: `Selamat datang, ${adminData.name}!`,
        timer: 2000,
        showConfirmButton: false,
      }).then(() => {
        navigate('/admin/dashboard');
      });
    } else {
      Swal.fire({
        icon: 'error',
        title: 'Login Gagal',
        text: 'Username atau password salah',
      });
    }
  };

  const handleDemoLogin = (role) => {
    const adminData = role === 'admin' 
      ? { id: 1, username: 'admin', name: 'Super Admin', role: 'super_admin' }
      : { id: 2, username: 'operator', name: 'Operator 1', role: 'operator' };
    
    adminLogin(adminData);
    Swal.fire({
      icon: 'success',
      title: 'Demo Login Berhasil',
      text: `Selamat datang, ${adminData.name}!`,
      timer: 2000,
      showConfirmButton: false,
    }).then(() => {
      navigate('/admin/dashboard');
    });
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary-dark via-primary-medium to-primary-light flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="bg-white/90 backdrop-blur-sm border border-white/20 shadow-2xl rounded-lg p-6 md:p-8">
          <div className="text-center mb-8">
            {/* Admin Logo */}
            <div className="mb-6">
              <img 
                src="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjgwIiB2aWV3Qm94PSIwIDAgMjAwIDgwIiBmaWxsPSJub25lIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8ZGVmcz4KPHN0eWxlPgouY2xzLTEgeyBmaWxsOiAjMzUzNzRCOyBmb250LWZhbWlseTogSW50ZXIsIHNhbnMtc2VyaWY7IGZvbnQtc2l6ZTogMjhweDsgZm9udC13ZWlnaHQ6IDcwMDsgfQouY2xzLTIgeyBmaWxsOiAjNTA3MjdCOyBmb250LWZhbWlseTogSW50ZXIsIHNhbnMtc2VyaWY7IGZvbnQtc2l6ZTogMTJweDsgZm9udC13ZWlnaHQ6IDQwMDsgfQouaWNvbiB7IGZpbGw6ICM3OEEwODM7IH0KPC9zdHlsZT4KPC9kZWZzPgo8IS0tIEdhbWUgSWNvbiAtLT4KPHJlY3QgeD0iMTAiIHk9IjE1IiB3aWR0aD0iMzAiIGhlaWdodD0iMjAiIHJ4PSIzIiBjbGFzcz0iaWNvbiIvPgo8Y2lyY2xlIGN4PSIxOCIgY3k9IjI4IiByPSIyIiBmaWxsPSJ3aGl0ZSIvPgo8Y2lyY2xlIGN4PSIzMiIgY3k9IjI4IiByPSIyIiBmaWxsPSJ3aGl0ZSIvPgo8cmVjdCB4PSIyMCIgeT0iMzgiIHdpZHRoPSI4IiBoZWlnaHQ9IjEwIiByeD0iMiIgY2xhc3M9Imljb24iLz4KPCEtLSBUZXh0IC0tPgo8dGV4dCB4PSI1NSIgeT0iNDAiIGNsYXNzPSJjbHMtMSI+V2F3IFN0b3JlPC90ZXh0Pgo8dGV4dCB4PSI1NSIgeT0iNTYiIGNsYXNzPSJjbHMtMiI+QWRtaW4gUGFuZWw8L3RleHQ+Cjwvc3ZnPgo="
                alt="Waw Store Admin"
                className="mx-auto h-16"
              />
            </div>
            <h1 className="text-xl md:text-2xl font-bold text-primary-dark mb-2">Admin Panel</h1>
            <p className="text-sm md:text-base text-text-secondary">Masuk ke panel administrasi</p>
          </div>
          
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-text-primary mb-2">Username</label>
              <input
                {...register('username')}
                className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                placeholder="Masukkan username"
              />
              {errors.username && <p className="text-danger text-xs mt-1">{errors.username.message}</p>}
            </div>
            
            <div>
              <label className="block text-sm font-medium text-text-primary mb-2">Password</label>
              <input
                type="password"
                {...register('password')}
                className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                placeholder="Masukkan password"
              />
              {errors.password && <p className="text-danger text-xs mt-1">{errors.password.message}</p>}
            </div>
            
            <button
              type="submit"
              className="w-full py-3 bg-primary-dark text-white rounded-lg hover:bg-primary-medium transition-colors focus:ring-2 focus:ring-primary-light focus:outline-none font-medium"
            >
              Masuk ke Admin
            </button>
          </form>
          
          <div className="mt-6 pt-6 border-t border-gray-200">
            <div className="text-center mb-4">
              <p className="text-sm text-text-secondary">Demo Login:</p>
            </div>
            
            <div className="space-y-2">
              <button
                onClick={() => handleDemoLogin('admin')}
                className="w-full py-2 bg-success text-white rounded-lg hover:bg-green-600 transition-colors focus:ring-2 focus:ring-green-500 focus:outline-none font-medium text-sm"
              >
                Demo Super Admin (admin/admin123)
              </button>
              
              <button
                onClick={() => handleDemoLogin('operator')}
                className="w-full py-2 bg-primary-light text-white rounded-lg hover:bg-primary-medium transition-colors focus:ring-2 focus:ring-blue-500 focus:outline-none font-medium text-sm"
              >
                Demo Operator (operator/operator123)
              </button>
            </div>
          </div>
          
          <div className="mt-4 text-center">
            <button
              onClick={() => navigate('/')}
              className="text-primary-light hover:text-primary-dark underline text-sm"
            >
              Kembali ke Beranda
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default AdminLogin;