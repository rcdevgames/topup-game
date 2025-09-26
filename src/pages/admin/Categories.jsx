import { useState } from 'react';
import { useAdminStore } from '../../store/adminStore';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import Swal from 'sweetalert2';

const categorySchema = yup.object({
  name: yup.string().required('Nama kategori wajib diisi'),
  description: yup.string().required('Deskripsi wajib diisi'),
  status: yup.string().required('Status wajib dipilih'),
});

function Categories() {
  const { categories, addCategory, updateCategory, deleteCategory } = useAdminStore();
  const [showForm, setShowForm] = useState(false);
  const [editingCategory, setEditingCategory] = useState(null);

  const { register, handleSubmit, formState: { errors }, reset } = useForm({
    resolver: yupResolver(categorySchema),
  });

  const onSubmit = (data) => {
    if (editingCategory) {
      updateCategory(editingCategory.id, data);
      Swal.fire({
        icon: 'success',
        title: 'Kategori Diperbarui',
        timer: 2000,
        showConfirmButton: false,
      });
      setEditingCategory(null);
    } else {
      addCategory(data);
      Swal.fire({
        icon: 'success',
        title: 'Kategori Ditambahkan',
        timer: 2000,
        showConfirmButton: false,
      });
    }
    setShowForm(false);
    reset();
  };

  const handleEdit = (category) => {
    setEditingCategory(category);
    setShowForm(true);
    reset(category);
  };

  const handleDelete = (id) => {
    Swal.fire({
      title: 'Hapus Kategori?',
      text: 'Kategori akan dihapus permanen',
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#EF4444',
      cancelButtonColor: '#6B7280',
      confirmButtonText: 'Ya, Hapus',
      cancelButtonText: 'Batal'
    }).then((result) => {
      if (result.isConfirmed) {
        deleteCategory(id);
        Swal.fire({
          icon: 'success',
          title: 'Kategori Dihapus',
          timer: 2000,
          showConfirmButton: false,
        });
      }
    });
  };

  const handleCancel = () => {
    setShowForm(false);
    setEditingCategory(null);
    reset();
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Kategori Produk</h1>
          <p className="text-gray-600">Kelola kategori produk game</p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          className="bg-primary-dark text-white px-4 py-2 rounded-lg hover:bg-primary-medium transition-colors flex items-center space-x-2"
        >
          <i className="fas fa-plus"></i>
          <span>Tambah Kategori</span>
        </button>
      </div>

      {/* Add/Edit Form */}
      {showForm && (
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">
            {editingCategory ? 'Edit Kategori' : 'Tambah Kategori'}
          </h3>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Nama Kategori</label>
                <input
                  {...register('name')}
                  className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                  placeholder="Masukkan nama kategori"
                />
                {errors.name && <p className="text-red-500 text-xs mt-1">{errors.name.message}</p>}
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
                <select
                  {...register('status')}
                  className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                >
                  <option value="">Pilih Status</option>
                  <option value="active">Aktif</option>
                  <option value="inactive">Nonaktif</option>
                </select>
                {errors.status && <p className="text-red-500 text-xs mt-1">{errors.status.message}</p>}
              </div>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Deskripsi</label>
              <textarea
                {...register('description')}
                rows={3}
                className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                placeholder="Masukkan deskripsi kategori"
              />
              {errors.description && <p className="text-red-500 text-xs mt-1">{errors.description.message}</p>}
            </div>
            
            <div className="flex space-x-3">
              <button
                type="submit"
                className="bg-success text-white px-4 py-2 rounded-lg hover:bg-green-600 transition-colors"
              >
                {editingCategory ? 'Update' : 'Simpan'}
              </button>
              <button
                type="button"
                onClick={handleCancel}
                className="bg-gray-500 text-white px-4 py-2 rounded-lg hover:bg-gray-600 transition-colors"
              >
                Batal
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Categories Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {categories.map((category) => (
          <div key={category.id} className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
            <div className="flex items-center justify-between mb-4">
              <div className="w-12 h-12 bg-primary-light bg-opacity-10 rounded-lg flex items-center justify-center">
                <i className="fas fa-gamepad text-primary-light text-xl"></i>
              </div>
              <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                category.status === 'active' 
                  ? 'bg-green-100 text-green-800' 
                  : 'bg-red-100 text-red-800'
              }`}>
                {category.status === 'active' ? 'Aktif' : 'Nonaktif'}
              </span>
            </div>
            
            <h3 className="text-lg font-semibold text-gray-900 mb-2">{category.name}</h3>
            <p className="text-gray-600 text-sm mb-4 line-clamp-2">{category.description}</p>
            
            <div className="flex items-center justify-between text-xs text-gray-500 mb-4">
              <span>Dibuat: {new Date(category.createdAt).toLocaleDateString('id-ID')}</span>
            </div>
            
            <div className="flex space-x-2">
              <button
                onClick={() => handleEdit(category)}
                className="flex-1 bg-primary-light text-white py-2 px-3 rounded-lg hover:bg-primary-medium transition-colors text-sm"
              >
                <i className="fas fa-edit mr-1"></i>
                Edit
              </button>
              <button
                onClick={() => handleDelete(category.id)}
                className="flex-1 bg-red-500 text-white py-2 px-3 rounded-lg hover:bg-red-600 transition-colors text-sm"
              >
                <i className="fas fa-trash mr-1"></i>
                Hapus
              </button>
            </div>
          </div>
        ))}
      </div>

      {categories.length === 0 && (
        <div className="bg-white rounded-lg shadow-md p-12 text-center">
          <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <i className="fas fa-tags text-gray-400 text-2xl"></i>
          </div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">Belum ada kategori</h3>
          <p className="text-gray-500 mb-4">Mulai dengan menambahkan kategori produk pertama</p>
          <button
            onClick={() => setShowForm(true)}
            className="bg-primary-dark text-white px-4 py-2 rounded-lg hover:bg-primary-medium transition-colors"
          >
            Tambah Kategori
          </button>
        </div>
      )}
    </div>
  );
}

export default Categories;