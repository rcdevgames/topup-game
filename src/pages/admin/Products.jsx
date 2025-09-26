import { useState } from 'react';
import { useAdminStore } from '../../store/adminStore';
import { useForm, useFieldArray } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import Swal from 'sweetalert2';

const productSchema = yup.object({
  name: yup.string().required('Nama produk wajib diisi'),
  categoryId: yup.number().required('Kategori wajib dipilih'),
  price: yup.string().required('Harga wajib diisi'),
  description: yup.string().required('Deskripsi wajib diisi'),
  image: yup.string().url('URL gambar tidak valid').required('URL gambar wajib diisi'),
  status: yup.string().required('Status wajib dipilih'),
  formConfig: yup.array().of(
    yup.object({
      field: yup.string().required(),
      label: yup.string().required(),
      type: yup.string().required(),
      required: yup.boolean(),
      placeholder: yup.string()
    })
  ).min(1, 'Minimal 1 field konfigurasi'),
});

function Products() {
  const { products, categories, addProduct, updateProduct, deleteProduct } = useAdminStore();
  const [showForm, setShowForm] = useState(false);
  const [editingProduct, setEditingProduct] = useState(null);

  const { register, handleSubmit, formState: { errors }, reset, watch, control } = useForm({
    resolver: yupResolver(productSchema),
    defaultValues: {
      formConfig: [{ field: 'gameAccount', label: 'ID Game', type: 'text', required: true, placeholder: 'Masukkan ID Game' }]
    }
  });

  const { fields, append, remove } = useFieldArray({
    control,
    name: 'formConfig'
  });

  const watchCategoryId = watch('categoryId');

  const onSubmit = (data) => {
    const categoryName = categories.find(c => c.id === parseInt(data.categoryId))?.name || '';
    const productData = { ...data, category: categoryName };
    
    if (editingProduct) {
      updateProduct(editingProduct.id, productData);
      Swal.fire({
        icon: 'success',
        title: 'Produk Diperbarui',
        timer: 2000,
        showConfirmButton: false,
      });
      setEditingProduct(null);
    } else {
      addProduct(productData);
      Swal.fire({
        icon: 'success',
        title: 'Produk Ditambahkan',
        timer: 2000,
        showConfirmButton: false,
      });
    }
    setShowForm(false);
    reset({
      formConfig: [{ field: 'gameAccount', label: 'ID Game', type: 'text', required: true, placeholder: 'Masukkan ID Game' }]
    });
  };

  const handleEdit = (product) => {
    setEditingProduct(product);
    setShowForm(true);
    reset(product);
  };

  const handleDelete = (id) => {
    Swal.fire({
      title: 'Hapus Produk?',
      text: 'Produk akan dihapus permanen',
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#EF4444',
      cancelButtonColor: '#6B7280',
      confirmButtonText: 'Ya, Hapus',
      cancelButtonText: 'Batal'
    }).then((result) => {
      if (result.isConfirmed) {
        deleteProduct(id);
        Swal.fire({
          icon: 'success',
          title: 'Produk Dihapus',
          timer: 2000,
          showConfirmButton: false,
        });
      }
    });
  };

  const handleCancel = () => {
    setShowForm(false);
    setEditingProduct(null);
    reset({
      formConfig: [{ field: 'gameAccount', label: 'ID Game', type: 'text', required: true, placeholder: 'Masukkan ID Game' }]
    });
  };

  const addFormField = () => {
    append({ field: '', label: '', type: 'text', required: false, placeholder: '' });
  };

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR'
    }).format(amount);
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Produk</h1>
          <p className="text-gray-600">Kelola produk topup game</p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          className="bg-primary-dark text-white px-4 py-2 rounded-lg hover:bg-primary-medium transition-colors flex items-center space-x-2"
        >
          <i className="fas fa-plus"></i>
          <span>Tambah Produk</span>
        </button>
      </div>

      {/* Add/Edit Form */}
      {showForm && (
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">
            {editingProduct ? 'Edit Produk' : 'Tambah Produk'}
          </h3>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            {/* Basic Info */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Nama Produk</label>
                <input
                  {...register('name')}
                  className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                  placeholder="Masukkan nama produk"
                />
                {errors.name && <p className="text-red-500 text-xs mt-1">{errors.name.message}</p>}
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Kategori</label>
                <select
                  {...register('categoryId', { valueAsNumber: true })}
                  className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                >
                  <option value="">Pilih Kategori</option>
                  {categories.map((category) => (
                    <option key={category.id} value={category.id}>{category.name}</option>
                  ))}
                </select>
                {errors.categoryId && <p className="text-red-500 text-xs mt-1">{errors.categoryId.message}</p>}
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Harga (Rp)</label>
                <input
                  {...register('price')}
                  type="number"
                  className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                  placeholder="Masukkan harga"
                />
                {errors.price && <p className="text-red-500 text-xs mt-1">{errors.price.message}</p>}
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
              <label className="block text-sm font-medium text-gray-700 mb-1">URL Gambar</label>
              <input
                {...register('image')}
                className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                placeholder="https://example.com/image.jpg"
              />
              {errors.image && <p className="text-red-500 text-xs mt-1">{errors.image.message}</p>}
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Deskripsi</label>
              <textarea
                {...register('description')}
                rows={3}
                className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                placeholder="Masukkan deskripsi produk"
              />
              {errors.description && <p className="text-red-500 text-xs mt-1">{errors.description.message}</p>}
            </div>

            {/* Dynamic Form Configuration */}
            <div className="border-t pt-6">
              <div className="flex items-center justify-between mb-4">
                <div>
                  <h4 className="text-lg font-medium text-gray-900">Konfigurasi Form Checkout</h4>
                  <p className="text-sm text-gray-600">Tentukan input fields yang akan muncul di halaman checkout</p>
                </div>
                <button
                  type="button"
                  onClick={addFormField}
                  className="bg-success text-white px-3 py-2 rounded-lg hover:bg-green-600 transition-colors text-sm flex items-center space-x-1"
                >
                  <i className="fas fa-plus text-xs"></i>
                  <span>Tambah Field</span>
                </button>
              </div>

              <div className="space-y-4">
                {fields.map((field, index) => (
                  <div key={field.id} className="p-4 border border-gray-200 rounded-lg">
                    <div className="grid grid-cols-1 md:grid-cols-5 gap-3">
                      <div>
                        <label className="block text-xs font-medium text-gray-700 mb-1">Field Name</label>
                        <input
                          {...register(`formConfig.${index}.field`)}
                          className="w-full p-2 border border-gray-300 rounded text-sm focus:ring-1 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                          placeholder="gameAccount"
                        />
                      </div>
                      
                      <div>
                        <label className="block text-xs font-medium text-gray-700 mb-1">Label</label>
                        <input
                          {...register(`formConfig.${index}.label`)}
                          className="w-full p-2 border border-gray-300 rounded text-sm focus:ring-1 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                          placeholder="ID Game"
                        />
                      </div>
                      
                      <div>
                        <label className="block text-xs font-medium text-gray-700 mb-1">Type</label>
                        <select
                          {...register(`formConfig.${index}.type`)}
                          className="w-full p-2 border border-gray-300 rounded text-sm focus:ring-1 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                        >
                          <option value="text">Text</option>
                          <option value="number">Number</option>
                          <option value="select">Select</option>
                        </select>
                      </div>
                      
                      <div>
                        <label className="block text-xs font-medium text-gray-700 mb-1">Placeholder</label>
                        <input
                          {...register(`formConfig.${index}.placeholder`)}
                          className="w-full p-2 border border-gray-300 rounded text-sm focus:ring-1 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                          placeholder="Masukkan ID Game"
                        />
                      </div>
                      
                      <div className="flex items-end space-x-2">
                        <label className="flex items-center">
                          <input
                            type="checkbox"
                            {...register(`formConfig.${index}.required`)}
                            className="mr-1"
                          />
                          <span className="text-xs text-gray-700">Required</span>
                        </label>
                        {fields.length > 1 && (
                          <button
                            type="button"
                            onClick={() => remove(index)}
                            className="p-2 text-red-500 hover:text-red-700"
                          >
                            <i className="fas fa-trash text-xs"></i>
                          </button>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
              {errors.formConfig && <p className="text-red-500 text-xs mt-1">{errors.formConfig.message}</p>}
            </div>
            
            <div className="flex space-x-3">
              <button
                type="submit"
                className="bg-success text-white px-4 py-2 rounded-lg hover:bg-green-600 transition-colors"
              >
                {editingProduct ? 'Update' : 'Simpan'}
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

      {/* Products Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {products.map((product) => (
          <div key={product.id} className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow">
            <img 
              src={product.image} 
              alt={product.name}
              className="w-full h-48 object-cover"
              onError={(e) => {
                e.target.src = 'https://images.unsplash.com/photo-1511512578047-dfb367046420?w=400&h=300&fit=crop&crop=center';
              }}
            />
            <div className="p-6">
              <div className="flex items-center justify-between mb-2">
                <span className="inline-flex px-2 py-1 text-xs font-semibold rounded-full bg-primary-light bg-opacity-10 text-primary-dark">
                  {product.category}
                </span>
                <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                  product.status === 'active' 
                    ? 'bg-green-100 text-green-800' 
                    : 'bg-red-100 text-red-800'
                }`}>
                  {product.status === 'active' ? 'Aktif' : 'Nonaktif'}
                </span>
              </div>
              
              <h3 className="text-lg font-semibold text-gray-900 mb-2">{product.name}</h3>
              <p className="text-gray-600 text-sm mb-3 line-clamp-2">{product.description}</p>
              <p className="text-xl font-bold text-primary-medium mb-4">{formatCurrency(parseInt(product.price))}</p>
              
              {/* Form Config Preview */}
              <div className="mb-4">
                <p className="text-xs text-gray-500 mb-2">Form Fields ({product.formConfig?.length || 0}):</p>
                <div className="flex flex-wrap gap-1">
                  {product.formConfig?.map((field, index) => (
                    <span key={index} className="inline-flex px-2 py-1 text-xs bg-gray-100 text-gray-600 rounded">
                      {field.label}
                    </span>
                  ))}
                </div>
              </div>
              
              <div className="flex space-x-2">
                <button
                  onClick={() => handleEdit(product)}
                  className="flex-1 bg-primary-light text-white py-2 px-3 rounded-lg hover:bg-primary-medium transition-colors text-sm"
                >
                  <i className="fas fa-edit mr-1"></i>
                  Edit
                </button>
                <button
                  onClick={() => handleDelete(product.id)}
                  className="flex-1 bg-red-500 text-white py-2 px-3 rounded-lg hover:bg-red-600 transition-colors text-sm"
                >
                  <i className="fas fa-trash mr-1"></i>
                  Hapus
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      {products.length === 0 && (
        <div className="bg-white rounded-lg shadow-md p-12 text-center">
          <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <i className="fas fa-box text-gray-400 text-2xl"></i>
          </div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">Belum ada produk</h3>
          <p className="text-gray-500 mb-4">Mulai dengan menambahkan produk pertama</p>
          <button
            onClick={() => setShowForm(true)}
            className="bg-primary-dark text-white px-4 py-2 rounded-lg hover:bg-primary-medium transition-colors"
          >
            Tambah Produk
          </button>
        </div>
      )}
    </div>
  );
}

export default Products;