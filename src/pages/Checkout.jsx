import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { useGameStore } from '../store/gameStore';
import { useAuthStore } from '../store/authStore';
import Swal from 'sweetalert2';

// Validation schema
const checkoutSchema = yup.object({
  gameAccount: yup.string().required('Akun game wajib dipilih atau diisi'),
  gameZone: yup.string().optional(),
  gameServer: yup.string().optional(),
  whatsapp: yup.string().matches(/^\d{10,13}$/, 'Nomor WhatsApp tidak valid (10-13 digit)').required('Nomor WhatsApp wajib diisi'),
  paymentMethod: yup.string().required('Metode pembayaran wajib dipilih'),
  voucherCode: yup.string().optional(),
});

function Checkout() {
  const { productId } = useParams();
  const navigate = useNavigate();
  const { products, gameAccounts, addTransaction } = useGameStore();
  const { user, isLoggedIn } = useAuthStore();
  
  const [voucherDiscount, setVoucherDiscount] = useState(0);
  const [voucherApplied, setVoucherApplied] = useState(false);
  const [isCheckingVoucher, setIsCheckingVoucher] = useState(false);
  const [selectedPaymentMethod, setSelectedPaymentMethod] = useState('');

  // Find product
  const product = products.find(p => p.id === parseInt(productId));

  const { register, handleSubmit, formState: { errors }, setValue, watch } = useForm({
    resolver: yupResolver(checkoutSchema),
    defaultValues: {
      gameAccount: '',
      gameZone: '',
      gameServer: '',
      whatsapp: isLoggedIn ? user.phone || '' : '',
      paymentMethod: '',
      voucherCode: '',
    }
  });

  const watchedVoucher = watch('voucherCode');

  useEffect(() => {
    if (isLoggedIn && user.phone) {
      setValue('whatsapp', user.phone);
    }
  }, [isLoggedIn, user.phone, setValue]);

  if (!product) {
    return (
      <div className="p-4">
        <div className="text-center py-12">
          <h2 className="text-xl font-semibold text-text-primary mb-2">Produk Tidak Ditemukan</h2>
          <p className="text-text-secondary mb-4">Produk yang Anda cari tidak tersedia</p>
          <button
            onClick={() => navigate('/')}
            className="px-6 py-2 bg-primary-light text-white rounded-lg hover:bg-primary-medium transition-colors"
          >
            Kembali ke Beranda
          </button>
        </div>
      </div>
    );
  }

  const paymentMethods = [
    { id: 'gopay', name: 'GoPay', fee: 0 },
    { id: 'ovo', name: 'OVO', fee: 0 },
    { id: 'dana', name: 'DANA', fee: 0 },
    { id: 'bca', name: 'BCA Virtual Account', fee: 2500 },
    { id: 'mandiri', name: 'Mandiri Virtual Account', fee: 2500 },
    { id: 'bni', name: 'BNI Virtual Account', fee: 2500 }
  ];

  const handleVoucherCheck = async () => {
    const voucherCode = watchedVoucher?.trim();
    if (!voucherCode) {
      Swal.fire({
        icon: 'warning',
        title: 'Kode Voucher Kosong',
        text: 'Silakan masukkan kode voucher terlebih dahulu',
      });
      return;
    }

    setIsCheckingVoucher(true);
    
    // Simulate API call
    setTimeout(() => {
      const validVouchers = {
        'NEWUSER10': { discount: 10, type: 'percentage' },
        'SAVE5K': { discount: 5000, type: 'fixed' },
        'WEEKEND20': { discount: 20, type: 'percentage' }
      };

      if (validVouchers[voucherCode.toUpperCase()]) {
        const voucher = validVouchers[voucherCode.toUpperCase()];
        let discount = 0;
        
        if (voucher.type === 'percentage') {
          discount = Math.floor((parseInt(product.price) * voucher.discount) / 100);
        } else {
          discount = voucher.discount;
        }
        
        setVoucherDiscount(discount);
        setVoucherApplied(true);
        Swal.fire({
          icon: 'success',
          title: 'Voucher Valid!',
          text: `Anda mendapat diskon Rp ${discount.toLocaleString('id-ID')}`,
          timer: 2000,
          showConfirmButton: false,
        });
      } else {
        setVoucherDiscount(0);
        setVoucherApplied(false);
        Swal.fire({
          icon: 'error',
          title: 'Voucher Tidak Valid',
          text: 'Kode voucher yang Anda masukkan tidak valid',
        });
      }
      setIsCheckingVoucher(false);
    }, 1500);
  };

  const calculateTotal = () => {
    const basePrice = parseInt(product.price);
    const selectedMethod = paymentMethods.find(m => m.id === selectedPaymentMethod);
    const fee = selectedMethod ? selectedMethod.fee : 0;
    return basePrice + fee - voucherDiscount;
  };

  const onSubmit = (data) => {
    // Create transaction
    const transaction = {
      id: `TRX${Date.now()}`,
      game: product.category,
      productName: product.name,
      amount: calculateTotal().toString(),
      originalAmount: product.price,
      fee: paymentMethods.find(m => m.id === data.paymentMethod)?.fee || 0,
      discount: voucherDiscount,
      voucherCode: voucherApplied ? data.voucherCode : null,
      status: 'pending',
      date: new Date().toISOString(),
      gameAccount: data.gameAccount,
      gameZone: data.gameZone || null,
      gameServer: data.gameServer || null,
      whatsapp: data.whatsapp,
      paymentMethod: data.paymentMethod,
      paymentMethodName: paymentMethods.find(m => m.id === data.paymentMethod)?.name
    };

    // Add to store
    addTransaction(transaction);

    // Show success and redirect
    Swal.fire({
      icon: 'success',
      title: 'Pesanan Berhasil!',
      text: 'Silakan lakukan pembayaran sesuai petunjuk',
      timer: 2000,
      showConfirmButton: false,
    }).then(() => {
      navigate(`/transactions/${transaction.id}`);
    });
  };

  return (
    <div className="p-4 pb-6 space-y-6">
      {/* Header */}
      <header className="flex items-center space-x-4 py-6 page-background">
        <button
          onClick={() => navigate('/')}
          className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          aria-label="Kembali"
        >
          <i className="fas fa-arrow-left text-xl text-text-primary"></i>
        </button>
        <div>
          <h1 className="text-xl font-bold text-primary-dark">Checkout</h1>
          <p className="text-sm text-text-secondary">Lengkapi data pembelian</p>
        </div>
      </header>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        {/* Product Info */}
        <div className="bg-white shadow-md rounded-lg p-4">
          <h3 className="font-semibold text-lg text-primary-dark mb-3">Produk yang Dibeli</h3>
          <div className="flex items-center space-x-4">
            <img 
              src={product.image} 
              alt={product.name}
              className="w-16 h-16 object-cover rounded-lg"
              onError={(e) => {
                e.target.src = 'https://images.unsplash.com/photo-1511512578047-dfb367046420?w=64&h=64&fit=crop&crop=center';
              }}
            />
            <div className="flex-1">
              <h4 className="font-medium">{product.name}</h4>
              <p className="text-sm text-text-secondary">{product.category}</p>
              <p className="font-bold text-primary-medium">Rp {parseInt(product.price).toLocaleString('id-ID')}</p>
            </div>
          </div>
        </div>

        {/* Game Account */}
        <div className="bg-white shadow-md rounded-lg p-4">
          <h3 className="font-semibold text-lg text-primary-dark mb-3">Data Akun Game</h3>
          
          {gameAccounts.length > 0 && isLoggedIn && (
            <div className="mb-4">
              <label className="block text-sm font-medium text-text-primary mb-2">Pilih Akun Tersimpan</label>
              <select
                {...register('gameAccount')}
                className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
              >
                <option value="">-- Pilih akun tersimpan atau isi manual --</option>
                {gameAccounts
                  .filter(acc => acc.game === product.category)
                  .map((account) => (
                    <option key={account.id} value={`${account.gameId} (${account.server})`}>
                      {account.game} - {account.gameId} ({account.server})
                    </option>
                  ))}
              </select>
            </div>
          )}
          
          {/* Dynamic Game Account Fields */}
          <div className="space-y-4">
            <p className="text-sm text-text-secondary mb-3">
              {gameAccounts.length > 0 && isLoggedIn ? 'Atau isi manual:' : 'Masukkan data akun game:'}
            </p>
            
            {/* Different input fields based on game type */}
            {product.category === 'Mobile Legends' && (
              <>
                <div>
                  <label className="block text-sm font-medium text-text-primary mb-2">User ID</label>
                  <input
                    {...register('gameAccount')}
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
                    placeholder="Masukkan User ID"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-text-primary mb-2">Zone ID</label>
                  <input
                    {...register('gameZone')}
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
                    placeholder="Masukkan Zone ID"
                  />
                </div>
              </>
            )}
            
            {product.category === 'Free Fire' && (
              <div>
                <label className="block text-sm font-medium text-text-primary mb-2">Player ID</label>
                <input
                  {...register('gameAccount')}
                  className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
                  placeholder="Masukkan Player ID"
                />
              </div>
            )}
            
            {product.category === 'PUBG Mobile' && (
              <div>
                <label className="block text-sm font-medium text-text-primary mb-2">Player ID</label>
                <input
                  {...register('gameAccount')}
                  className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
                  placeholder="Masukkan Player ID"
                />
              </div>
            )}
            
            {product.category === 'Genshin Impact' && (
              <>
                <div>
                  <label className="block text-sm font-medium text-text-primary mb-2">UID</label>
                  <input
                    {...register('gameAccount')}
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
                    placeholder="Masukkan UID"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-text-primary mb-2">Server</label>
                  <select
                    {...register('gameServer')}
                    className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
                  >
                    <option value="">Pilih Server</option>
                    <option value="America">America</option>
                    <option value="Europe">Europe</option>
                    <option value="Asia">Asia</option>
                    <option value="TW_HK_MO">TW, HK, MO</option>
                  </select>
                </div>
              </>
            )}
            
            {/* Default fallback for other games */}
            {!['Mobile Legends', 'Free Fire', 'PUBG Mobile', 'Genshin Impact'].includes(product.category) && (
              <div>
                <label className="block text-sm font-medium text-text-primary mb-2">ID Game & Server</label>
                <input
                  {...register('gameAccount')}
                  className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
                  placeholder="Contoh: 123456789 (2001) atau 123456789,2001"
                />
              </div>
            )}
            
            {errors.gameAccount && <p className="text-danger text-xs mt-1">{errors.gameAccount.message}</p>}
          </div>
        </div>

        {/* Voucher */}
        <div className="bg-white shadow-md rounded-lg p-4">
          <h3 className="font-semibold text-lg text-primary-dark mb-3">Kode Voucher</h3>
          
          {/* Desktop Layout (side by side) */}
          <div className="hidden md:flex space-x-2">
            <input
              {...register('voucherCode')}
              className="flex-1 p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
              placeholder="Masukkan kode voucher (opsional)"
            />
            <button
              type="button"
              onClick={voucherApplied ? () => {
                setVoucherApplied(false);
                setVoucherDiscount(0);
                setValue('voucherCode', '');
              } : handleVoucherCheck}
              disabled={isCheckingVoucher}
              className={`px-4 py-3 text-white rounded-lg transition-colors focus:ring-2 focus:outline-none disabled:opacity-50 ${
                voucherApplied 
                  ? 'bg-red-500 hover:bg-red-600 focus:ring-red-500' 
                  : 'bg-primary-light hover:bg-primary-medium focus:ring-blue-500'
              }`}
            >
              {isCheckingVoucher ? 'Checking...' : voucherApplied ? 'Hapus' : 'Check'}
            </button>
          </div>
          
          {/* Mobile Layout (stacked) */}
          <div className="md:hidden space-y-2">
            <input
              {...register('voucherCode')}
              className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
              placeholder="Masukkan kode voucher (opsional)"
            />
            <button
              type="button"
              onClick={voucherApplied ? () => {
                setVoucherApplied(false);
                setVoucherDiscount(0);
                setValue('voucherCode', '');
              } : handleVoucherCheck}
              disabled={isCheckingVoucher}
              className={`w-full py-3 text-white rounded-lg transition-colors focus:ring-2 focus:outline-none disabled:opacity-50 ${
                voucherApplied 
                  ? 'bg-red-500 hover:bg-red-600 focus:ring-red-500' 
                  : 'bg-primary-light hover:bg-primary-medium focus:ring-blue-500'
              }`}
            >
              {isCheckingVoucher ? 'Checking...' : voucherApplied ? 'Hapus Voucher' : 'Cek Voucher'}
            </button>
          </div>
          
          {voucherApplied && (
            <div className="mt-3 p-3 bg-green-50 border border-green-200 rounded-lg flex items-center space-x-2">
              <i className="fas fa-check-circle text-success"></i>
              <span className="text-green-800 font-medium">
                Voucher berhasil diterapkan! Diskon: Rp {voucherDiscount.toLocaleString('id-ID')}
              </span>
            </div>
          )}
          
          <div className="mt-2 text-xs text-text-secondary">
            <p>Kode voucher demo: <span className="font-mono">NEWUSER10</span>, <span className="font-mono">SAVE5K</span>, <span className="font-mono">WEEKEND20</span></p>
          </div>
        </div>

        {/* Payment Method */}
        <div className="bg-white shadow-md rounded-lg p-4">
          <h3 className="font-semibold text-lg text-primary-dark mb-3">Metode Pembayaran</h3>
          <div className="space-y-2">
            {paymentMethods.map((method) => (
              <label key={method.id} className="flex items-center p-3 border border-gray-300 rounded-lg hover:bg-gray-50 cursor-pointer">
                <input
                  type="radio"
                  {...register('paymentMethod')}
                  value={method.id}
                  onChange={(e) => setSelectedPaymentMethod(e.target.value)}
                  className="mr-3"
                />
                <div className="flex-1">
                  <span className="font-medium">{method.name}</span>
                  {method.fee > 0 && (
                    <span className="text-sm text-text-secondary ml-2">
                      (+Rp {method.fee.toLocaleString('id-ID')} biaya admin)
                    </span>
                  )}
                </div>
              </label>
            ))}
          </div>
          {errors.paymentMethod && <p className="text-danger text-xs mt-1">{errors.paymentMethod.message}</p>}
        </div>

        {/* WhatsApp */}
        <div className="bg-white shadow-md rounded-lg p-4">
          <h3 className="font-semibold text-lg text-primary-dark mb-3">Nomor WhatsApp</h3>
          <input
            {...register('whatsapp')}
            className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
            placeholder="08XXXXXXXXXX"
          />
          {errors.whatsapp && <p className="text-danger text-xs mt-1">{errors.whatsapp.message}</p>}
          <p className="text-xs text-text-secondary mt-1">Notifikasi status pembayaran akan dikirim ke nomor ini</p>
        </div>

        {/* Order Summary */}
        <div className="bg-white shadow-md rounded-lg p-4">
          <h3 className="font-semibold text-lg text-primary-dark mb-3">Ringkasan Pesanan</h3>
          <div className="space-y-2">
            <div className="flex justify-between">
              <span>Harga Produk:</span>
              <span>Rp {parseInt(product.price).toLocaleString('id-ID')}</span>
            </div>
            {selectedPaymentMethod && paymentMethods.find(m => m.id === selectedPaymentMethod)?.fee > 0 && (
              <div className="flex justify-between">
                <span>Biaya Admin:</span>
                <span>Rp {paymentMethods.find(m => m.id === selectedPaymentMethod).fee.toLocaleString('id-ID')}</span>
              </div>
            )}
            {voucherApplied && (
              <div className="flex justify-between text-success">
                <span>Diskon Voucher:</span>
                <span>-Rp {voucherDiscount.toLocaleString('id-ID')}</span>
              </div>
            )}
            <hr />
            <div className="flex justify-between text-lg font-bold">
              <span>Total Pembayaran:</span>
              <span className="text-primary-medium">Rp {calculateTotal().toLocaleString('id-ID')}</span>
            </div>
          </div>
        </div>

        {/* Submit Button */}
        <button
          type="submit"
          className="w-full py-4 bg-success text-white rounded-lg hover:bg-green-600 transition-colors focus:ring-2 focus:ring-green-500 focus:outline-none font-medium text-lg"
        >
          Bayar Sekarang
        </button>
      </form>
    </div>
  );
}

export default Checkout;