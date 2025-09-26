import { useState } from 'react';
import { useAdminStore } from '../../store/adminStore';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import Swal from 'sweetalert2';

const voucherSchema = yup.object({
  code: yup.string().min(3, 'Kode minimal 3 karakter').required('Kode voucher wajib diisi'),
  type: yup.string().required('Tipe voucher wajib dipilih'),
  value: yup.number().positive('Nilai harus positif').required('Nilai wajib diisi'),
  description: yup.string().required('Deskripsi wajib diisi'),
  applicationType: yup.string().required('Tipe aplikasi wajib dipilih'),
  applicableIds: yup.array().when('applicationType', {
    is: (val) => val === 'category' || val === 'product',
    then: (schema) => schema.min(1, 'Minimal pilih 1 item'),
    otherwise: (schema) => schema
  }),
  quota: yup.number().positive('Quota harus positif').required('Quota wajib diisi'),
  startDate: yup.date().required('Tanggal mulai wajib diisi'),
  endDate: yup.date().min(yup.ref('startDate'), 'Tanggal berakhir harus setelah tanggal mulai').required('Tanggal berakhir wajib diisi'),
  status: yup.string().required('Status wajib dipilih'),
});

function Vouchers() {
  const { vouchers, categories, products, addVoucher, updateVoucher, deleteVoucher } = useAdminStore();
  const [showForm, setShowForm] = useState(false);
  const [editingVoucher, setEditingVoucher] = useState(null);

  const { register, handleSubmit, formState: { errors }, reset, watch, setValue } = useForm({
    resolver: yupResolver(voucherSchema),
  });

  const watchApplicationType = watch('applicationType');
  const watchType = watch('type');

  const onSubmit = (data) => {
    // Convert applicableIds to numbers
    if (data.applicableIds) {
      data.applicableIds = data.applicableIds.map(id => parseInt(id));
    }
    
    if (editingVoucher) {
      updateVoucher(editingVoucher.id, data);
      Swal.fire({
        icon: 'success',
        title: 'Voucher Diperbarui',
        timer: 2000,
        showConfirmButton: false,
      });
      setEditingVoucher(null);
    } else {
      addVoucher(data);
      Swal.fire({
        icon: 'success',
        title: 'Voucher Ditambahkan',
        timer: 2000,
        showConfirmButton: false,
      });
    }
    setShowForm(false);
    reset();
  };

  const handleEdit = (voucher) => {
    setEditingVoucher(voucher);
    setShowForm(true);
    reset({
      ...voucher,
      startDate: voucher.startDate,
      endDate: voucher.endDate
    });
  };

  const handleDelete = (id) => {
    Swal.fire({
      title: 'Hapus Voucher?',
      text: 'Voucher akan dihapus permanen',
      icon: 'warning',
      showCancelButton: true,
      confirmButtonColor: '#EF4444',
      cancelButtonColor: '#6B7280',
      confirmButtonText: 'Ya, Hapus',
      cancelButtonText: 'Batal'
    }).then((result) => {
      if (result.isConfirmed) {
        deleteVoucher(id);
        Swal.fire({
          icon: 'success',
          title: 'Voucher Dihapus',
          timer: 2000,
          showConfirmButton: false,
        });
      }
    });
  };

  const handleCancel = () => {
    setShowForm(false);
    setEditingVoucher(null);
    reset();
  };

  const generateVoucherCode = () => {
    const codes = ['SAVE', 'DISC', 'PROMO', 'SPECIAL', 'BONUS'];
    const numbers = Math.floor(Math.random() * 90) + 10;
    const randomCode = codes[Math.floor(Math.random() * codes.length)] + numbers;
    setValue('code', randomCode);
  };

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR'
    }).format(amount);
  };

  const getApplicableNames = (voucher) => {
    if (voucher.applicationType === 'all') return 'Semua Produk';
    if (voucher.applicationType === 'category') {
      return voucher.applicableIds.map(id => 
        categories.find(c => c.id === id)?.name
      ).filter(Boolean).join(', ');
    }
    if (voucher.applicationType === 'product') {
      return voucher.applicableIds.map(id => 
        products.find(p => p.id === id)?.name
      ).filter(Boolean).join(', ');
    }
    return '';
  };

  const isVoucherExpired = (endDate) => {
    return new Date(endDate) < new Date();
  };

  const getVoucherStatus = (voucher) => {
    if (voucher.status === 'inactive') return { text: 'Nonaktif', color: 'bg-gray-100 text-gray-800' };
    if (isVoucherExpired(voucher.endDate)) return { text: 'Expired', color: 'bg-red-100 text-red-800' };
    if (voucher.usedCount >= voucher.quota) return { text: 'Habis', color: 'bg-orange-100 text-orange-800' };
    return { text: 'Aktif', color: 'bg-green-100 text-green-800' };
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Voucher</h1>
          <p className="text-gray-600">Kelola voucher diskon dan promo</p>
        </div>
        <button
          onClick={() => setShowForm(true)}
          className="bg-primary-dark text-white px-4 py-2 rounded-lg hover:bg-primary-medium transition-colors flex items-center space-x-2"
        >
          <i className="fas fa-plus"></i>
          <span>Tambah Voucher</span>
        </button>
      </div>

      {/* Add/Edit Form */}
      {showForm && (
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-lg font-semibold text-gray-900 mb-4">
            {editingVoucher ? 'Edit Voucher' : 'Tambah Voucher'}
          </h3>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            {/* Basic Info */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Kode Voucher</label>
                <div className="flex space-x-2">
                  <input
                    {...register('code')}
                    className="flex-1 p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none uppercase"
                    placeholder="SAVE50"
                  />
                  <button
                    type="button"
                    onClick={generateVoucherCode}
                    className="px-3 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition-colors text-sm"
                  >
                    <i className="fas fa-random"></i>
                  </button>
                </div>
                {errors.code && <p className="text-red-500 text-xs mt-1">{errors.code.message}</p>}
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Tipe Voucher</label>
                <select
                  {...register('type')}
                  className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                >
                  <option value="">Pilih Tipe</option>
                  <option value="percentage">Persentase (%)</option>
                  <option value="fixed">Nominal (Rp)</option>
                </select>
                {errors.type && <p className="text-red-500 text-xs mt-1">{errors.type.message}</p>}
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Nilai {watchType === 'percentage' ? '(%)' : '(Rp)'}
                </label>
                <input
                  {...register('value', { valueAsNumber: true })}
                  type="number"
                  className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                  placeholder={watchType === 'percentage' ? '10' : '5000'}
                />
                {errors.value && <p className="text-red-500 text-xs mt-1">{errors.value.message}</p>}
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Quota Penggunaan</label>
                <input
                  {...register('quota', { valueAsNumber: true })}
                  type="number"
                  className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                  placeholder="100"
                />
                {errors.quota && <p className="text-red-500 text-xs mt-1">{errors.quota.message}</p>}
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Deskripsi</label>
              <textarea
                {...register('description')}
                rows={3}
                className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                placeholder="Deskripsi voucher"
              />
              {errors.description && <p className="text-red-500 text-xs mt-1">{errors.description.message}</p>}
            </div>

            {/* Date Range */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Tanggal Mulai</label>
                <input
                  {...register('startDate')}
                  type="date"
                  className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                />
                {errors.startDate && <p className="text-red-500 text-xs mt-1">{errors.startDate.message}</p>}
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Tanggal Berakhir</label>
                <input
                  {...register('endDate')}
                  type="date"
                  className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                />
                {errors.endDate && <p className="text-red-500 text-xs mt-1">{errors.endDate.message}</p>}
              </div>
            </div>

            {/* Applicable Products/Categories */}
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Berlaku Untuk</label>
                <select
                  {...register('applicationType')}
                  className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none"
                >
                  <option value="">Pilih Aplikasi</option>
                  <option value="all">Semua Produk</option>
                  <option value="category">Kategori Tertentu</option>
                  <option value="product">Produk Tertentu</option>
                </select>
                {errors.applicationType && <p className="text-red-500 text-xs mt-1">{errors.applicationType.message}</p>}
              </div>

              {watchApplicationType === 'category' && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Pilih Kategori</label>
                  <select
                    {...register('applicableIds')}
                    multiple
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none h-32"
                  >
                    {categories.map((category) => (
                      <option key={category.id} value={category.id}>{category.name}</option>
                    ))}
                  </select>
                  <p className="text-xs text-gray-500 mt-1">Tahan Ctrl/Cmd untuk memilih beberapa kategori</p>
                  {errors.applicableIds && <p className="text-red-500 text-xs mt-1">{errors.applicableIds.message}</p>}
                </div>
              )}

              {watchApplicationType === 'product' && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Pilih Produk</label>
                  <select
                    {...register('applicableIds')}
                    multiple
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-light focus:border-primary-light focus:outline-none h-32"
                  >
                    {products.map((product) => (
                      <option key={product.id} value={product.id}>{product.name} - {product.category}</option>
                    ))}
                  </select>
                  <p className="text-xs text-gray-500 mt-1">Tahan Ctrl/Cmd untuk memilih beberapa produk</p>
                  {errors.applicableIds && <p className="text-red-500 text-xs mt-1">{errors.applicableIds.message}</p>}
                </div>
              )}
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
            
            <div className="flex space-x-3">
              <button
                type="submit"
                className="bg-success text-white px-4 py-2 rounded-lg hover:bg-green-600 transition-colors"
              >
                {editingVoucher ? 'Update' : 'Simpan'}
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

      {/* Vouchers Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {vouchers.map((voucher) => {
          const status = getVoucherStatus(voucher);
          return (
            <div key={voucher.id} className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow">
              <div className="flex items-center justify-between mb-4">
                <div className="flex items-center space-x-2">
                  <div className="w-10 h-10 bg-yellow-500 bg-opacity-10 rounded-lg flex items-center justify-center">
                    <i className="fas fa-ticket-alt text-yellow-600"></i>
                  </div>
                  <div>
                    <h3 className="font-bold text-lg text-gray-900">{voucher.code}</h3>
                    <p className="text-xs text-gray-500">
                      {voucher.type === 'percentage' ? `${voucher.value}%` : formatCurrency(voucher.value)}
                    </p>
                  </div>
                </div>
                <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${status.color}`}>
                  {status.text}
                </span>
              </div>
              
              <p className="text-gray-600 text-sm mb-4 line-clamp-2">{voucher.description}</p>
              
              <div className="space-y-2 text-xs text-gray-500 mb-4">
                <div className="flex justify-between">
                  <span>Berlaku:</span>
                  <span>{getApplicableNames(voucher)}</span>
                </div>
                <div className="flex justify-between">
                  <span>Penggunaan:</span>
                  <span>{voucher.usedCount}/{voucher.quota}</span>
                </div>
                <div className="flex justify-between">
                  <span>Periode:</span>
                  <span>{new Date(voucher.startDate).toLocaleDateString('id-ID')} - {new Date(voucher.endDate).toLocaleDateString('id-ID')}</span>
                </div>
              </div>

              {/* Usage Bar */}
              <div className="mb-4">
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div 
                    className="bg-primary-light h-2 rounded-full transition-all duration-300"
                    style={{ width: `${(voucher.usedCount / voucher.quota) * 100}%` }}
                  ></div>
                </div>
              </div>
              
              <div className="flex space-x-2">
                <button
                  onClick={() => handleEdit(voucher)}
                  className="flex-1 bg-primary-light text-white py-2 px-3 rounded-lg hover:bg-primary-medium transition-colors text-sm"
                >
                  <i className="fas fa-edit mr-1"></i>
                  Edit
                </button>
                <button
                  onClick={() => handleDelete(voucher.id)}
                  className="flex-1 bg-red-500 text-white py-2 px-3 rounded-lg hover:bg-red-600 transition-colors text-sm"
                >
                  <i className="fas fa-trash mr-1"></i>
                  Hapus
                </button>
              </div>
            </div>
          );
        })}
      </div>

      {vouchers.length === 0 && (
        <div className="bg-white rounded-lg shadow-md p-12 text-center">
          <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <i className="fas fa-ticket-alt text-gray-400 text-2xl"></i>
          </div>
          <h3 className="text-lg font-medium text-gray-900 mb-2">Belum ada voucher</h3>
          <p className="text-gray-500 mb-4">Mulai dengan membuat voucher pertama</p>
          <button
            onClick={() => setShowForm(true)}
            className="bg-primary-dark text-white px-4 py-2 rounded-lg hover:bg-primary-medium transition-colors"
          >
            Tambah Voucher
          </button>
        </div>
      )}
    </div>
  );
}

export default Vouchers;