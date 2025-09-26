import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import { loginSchema } from '../utils/validation';
import Swal from 'sweetalert2';

function Login() {
  const navigate = useNavigate();
  const { login } = useAuthStore();
  
  const { register, handleSubmit, formState: { errors } } = useForm({
    resolver: yupResolver(loginSchema),
  });

  const onSubmit = (data) => {
    // Simple demo login - in real app, this would call an API
    if (data.phone === '08123456789' && data.password === 'password') {
      const userData = {
        name: 'Demo User',
        phone: data.phone,
        password: data.password
      };
      
      login(userData);
      Swal.fire({
        icon: 'success',
        title: 'Login Berhasil',
        text: 'Selamat datang!',
        timer: 2000,
        showConfirmButton: false,
      }).then(() => {
        navigate('/');
      });
    } else {
      Swal.fire({
        icon: 'error',
        title: 'Login Gagal',
        text: 'Nomor HP atau password salah',
      });
    }
  };

  const handleDemoLogin = () => {
    const userData = {
      name: 'Demo User',
      phone: '08123456789',
      password: 'password'
    };
    
    login(userData);
    Swal.fire({
      icon: 'success',
      title: 'Demo Login Berhasil',
      text: 'Selamat datang!',
      timer: 2000,
      showConfirmButton: false,
    }).then(() => {
      navigate('/');
    });
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="bg-white shadow-lg rounded-lg p-6 page-background">
          <div className="text-center mb-6">
            {/* Logo Waw Store */}
            <div className="mb-6">
              <img 
                src="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjgwIiB2aWV3Qm94PSIwIDAgMjAwIDgwIiBmaWxsPSJub25lIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8ZGVmcz4KPHN0eWxlPgouY2xzLTEgeyBmaWxsOiAjMzUzNzRCOyBmb250LWZhbWlseTogSW50ZXIsIHNhbnMtc2VyaWY7IGZvbnQtc2l6ZTogMjhweDsgZm9udC13ZWlnaHQ6IDcwMDsgfQouY2xzLTIgeyBmaWxsOiAjNTA3MjdCOyBmb250LWZhbWlseTogSW50ZXIsIHNhbnMtc2VyaWY7IGZvbnQtc2l6ZTogMTJweDsgZm9udC13ZWlnaHQ6IDQwMDsgfQouaWNvbiB7IGZpbGw6ICM3OEEwODM7IH0KPC9zdHlsZT4KPC9kZWZzPgo8IS0tIEdhbWUgSWNvbiAtLT4KPHJlY3QgeD0iMTAiIHk9IjE1IiB3aWR0aD0iMzAiIGhlaWdodD0iMjAiIHJ4PSIzIiBjbGFzcz0iaWNvbiIvPgo8Y2lyY2xlIGN4PSIxOCIgY3k9IjI4IiByPSIyIiBmaWxsPSJ3aGl0ZSIvPgo8Y2lyY2xlIGN4PSIzMiIgY3k9IjI4IiByPSIyIiBmaWxsPSJ3aGl0ZSIvPgo8cmVjdCB4PSIyMCIgeT0iMzgiIHdpZHRoPSI4IiBoZWlnaHQ9IjEwIiByeD0iMiIgY2xhc3M9Imljb24iLz4KPCEtLSBUZXh0IC0tPgo8dGV4dCB4PSI1NSIgeT0iNDAiIGNsYXNzPSJjbHMtMSI+V2F3IFN0b3JlPC90ZXh0Pgo8dGV4dCB4PSI1NSIgeT0iNTYiIGNsYXNzPSJjbHMtMiI+VG9wdXAgR2FtZSBPbmxpbmU8L3RleHQ+Cjwvc3ZnPgo="
                alt="Waw Store Logo"
                className="mx-auto h-16"
              />
            </div>
            <h1 className="text-2xl font-bold text-primary-dark mb-2">Masuk</h1>
            <p className="text-text-secondary">Masuk ke akun topup game kamu</p>
          </div>
          
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-text-primary mb-2">Nomor HP</label>
              <input
                {...register('phone')}
                className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
                placeholder="08XXXXXXXXXX"
              />
              {errors.phone && <p className="text-danger text-xs mt-1">{errors.phone.message}</p>}
            </div>
            
            <div>
              <label className="block text-sm font-medium text-text-primary mb-2">Password</label>
              <input
                type="password"
                {...register('password')}
                className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
                placeholder="Masukkan password"
              />
              {errors.password && <p className="text-danger text-xs mt-1">{errors.password.message}</p>}
            </div>
            
            <button
              type="submit"
              className="w-full py-3 bg-primary-light text-white rounded-lg hover:bg-primary-medium transition-colors focus:ring-2 focus:ring-blue-500 focus:outline-none font-medium"
            >
              Masuk
            </button>
          </form>
          
          <div className="mt-6 pt-6 border-t border-gray-200">
            <div className="text-center mb-4">
              <p className="text-sm text-text-secondary">Untuk demo, gunakan:</p>
              <p className="text-xs text-text-secondary mt-1">HP: 08123456789 | Password: password</p>
            </div>
            
            <button
              onClick={handleDemoLogin}
              className="w-full py-3 bg-success text-white rounded-lg hover:bg-green-600 transition-colors focus:ring-2 focus:ring-green-500 focus:outline-none font-medium"
            >
              Demo Login
            </button>
          </div>
          
          <div className="mt-4 text-center">
            <button
              onClick={() => navigate('/')}
              className="text-primary-light hover:text-primary-medium underline text-sm"
            >
              Kembali ke Beranda
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Login;