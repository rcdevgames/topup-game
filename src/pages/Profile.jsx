import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { useAuthStore } from '../store/authStore';
import { useGameStore } from '../store/gameStore';
import { profileSchema } from '../utils/validation';
import GameAccountForm from '../components/GameAccountForm';
import Swal from 'sweetalert2';

function Profile() {
  const { user, updateUser, logout } = useAuthStore();
  const { gameAccounts, addGameAccount, updateGameAccount, deleteGameAccount } = useGameStore();
  const [editingAccount, setEditingAccount] = useState(null);
  const [showAddForm, setShowAddForm] = useState(false);
  
  const { register, handleSubmit, formState: { errors } } = useForm({
    resolver: yupResolver(profileSchema),
    defaultValues: {
      name: user.name || '',
      phone: user.phone || '',
      password: ''
    },
  });

  const onSubmit = (data) => {
    // Remove password if empty
    const updateData = { ...data };
    if (!updateData.password) {
      delete updateData.password;
    }
    
    updateUser(updateData);
    Swal.fire({
      icon: 'success',
      title: 'Profil Diperbarui',
      text: 'Data profil berhasil disimpan',
      timer: 2000,
      showConfirmButton: false,
    });
  };

  const handleAddAccount = (data) => {
    addGameAccount(data);
    setShowAddForm(false);
    Swal.fire({
      icon: 'success',
      title: 'Akun Game Ditambahkan',
      timer: 2000,
      showConfirmButton: false,
    });
  };

  const handleEditAccount = (data) => {
    updateGameAccount(editingAccount.id, { ...data, id: editingAccount.id });
    setEditingAccount(null);
    Swal.fire({
      icon: 'success',
      title: 'Akun Game Diperbarui',
      timer: 2000,
      showConfirmButton: false,
    });
  };

  const handleDeleteAccount = (id) => {
    Swal.fire({
      title: 'Hapus Akun Game?',
      text: 'Akun game akan dihapus permanen',
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#EF4444',
      cancelButtonColor: '#6B7280',
      confirmButtonText: 'Ya, Hapus',
      cancelButtonText: 'Batal'
    }).then((result) => {
      if (result.isConfirmed) {
        deleteGameAccount(id);
        Swal.fire({
          icon: 'success',
          title: 'Akun Game Dihapus',
          timer: 2000,
          showConfirmButton: false,
        });
      }
    });
  };

  const handleLogout = () => {
    Swal.fire({
      title: 'Keluar dari Akun?',
      icon: 'question',
      showCancelButton: true,
      confirmButtonColor: '#EF4444',
      cancelButtonColor: '#6B7280',
      confirmButtonText: 'Ya, Keluar',
      cancelButtonText: 'Batal'
    }).then((result) => {
      if (result.isConfirmed) {
        logout();
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
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-green-50">
      <div className="p-4 pb-6 space-y-6">
        <header className="py-6 page-background">
          <h1 className="text-2xl font-bold text-primary-dark mb-2">Profil</h1>
          <p className="text-text-secondary">Kelola informasi akun dan game kamu</p>
        </header>
      
      {/* User Info Form */}
      <section>
        <form onSubmit={handleSubmit(onSubmit)} className="bg-white/80 backdrop-blur-sm border border-white/20 shadow-lg rounded-lg p-4 space-y-4">
          <h2 className="text-lg font-semibold text-primary-dark mb-4">Informasi Pribadi</h2>
          
          <div>
            <label className="block text-sm font-medium text-text-primary mb-2">Nama Lengkap</label>
            <input
              {...register('name')}
              className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
              placeholder="Masukkan nama lengkap"
            />
            {errors.name && <p className="text-danger text-xs mt-1">{errors.name.message}</p>}
          </div>
          
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
            <label className="block text-sm font-medium text-text-primary mb-2">Ganti Password</label>
            <input
              type="password"
              {...register('password')}
              className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
              placeholder="Kosongkan jika tidak ingin mengubah"
            />
            {errors.password && <p className="text-danger text-xs mt-1">{errors.password.message}</p>}
          </div>
          
          <button
            type="submit"
            className="w-full py-3 bg-primary-light text-white rounded-lg hover:bg-primary-medium transition-colors focus:ring-2 focus:ring-blue-500 focus:outline-none font-medium"
          >
            Simpan Perubahan
          </button>
        </form>
      </section>
      
      {/* Game Accounts CRUD */}
      <section>
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-lg font-semibold text-primary-dark">Akun Game ({gameAccounts.length})</h2>
          <button
            onClick={() => setShowAddForm(!showAddForm)}
            className="px-4 py-2 bg-success text-white rounded-lg hover:bg-green-600 transition-colors focus:ring-2 focus:ring-green-500 focus:outline-none"
          >
            {showAddForm ? 'Batal' : 'Tambah Akun'}
          </button>
        </div>
        
        {/* Add Form */}
        {showAddForm && (
          <div className="mb-4">
            <GameAccountForm onSubmit={handleAddAccount} />
          </div>
        )}
        
        {/* Edit Form */}
        {editingAccount && (
          <div className="mb-4">
            <GameAccountForm 
              onSubmit={handleEditAccount} 
              initialData={editingAccount}
              isEdit={true}
            />
            <button
              onClick={() => setEditingAccount(null)}
              className="mt-2 px-4 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition-colors"
            >
              Batal Edit
            </button>
          </div>
        )}
        
        {/* Game Accounts List */}
        {gameAccounts.length === 0 ? (
          <div className="bg-white/80 backdrop-blur-sm border border-white/20 shadow-lg rounded-lg p-8 text-center">
            <div className="mb-4">
              <i className="fas fa-gamepad text-5xl text-text-secondary"></i>
            </div>
            <p className="text-text-secondary">Belum ada akun game tersimpan</p>
          </div>
        ) : (
          <div className="space-y-3">
            {gameAccounts.map((account) => (
              <div key={account.id} className="bg-white/80 backdrop-blur-sm border border-white/20 shadow-lg rounded-lg p-4">
                <div className="flex justify-between items-start">
                  <div className="flex-1">
                    <h3 className="font-semibold text-primary-dark">{account.game}</h3>
                    <p className="text-sm text-text-secondary">ID: {account.gameId}</p>
                    <p className="text-sm text-text-secondary">Server: {account.server}</p>
                  </div>
                  <div className="flex space-x-2 ml-4">
                    <button
                      onClick={() => setEditingAccount(account)}
                      className="px-3 py-1 bg-primary-light text-white rounded text-sm hover:bg-primary-medium transition-colors"
                    >
                      Edit
                    </button>
                    <button
                      onClick={() => handleDeleteAccount(account.id)}
                      className="px-3 py-1 bg-danger text-white rounded text-sm hover:bg-red-600 transition-colors"
                    >
                      Hapus
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </section>
      
        {/* Logout Button */}
        <section>
          <button
            onClick={handleLogout}
            className="w-full py-3 bg-danger text-white rounded-lg hover:bg-red-600 transition-colors focus:ring-2 focus:ring-red-500 focus:outline-none font-medium"
          >
            Keluar dari Akun
          </button>
        </section>
      </div>
    </div>
  );
}

export default Profile;