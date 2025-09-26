import { useParams, useNavigate } from 'react-router-dom';
import { useGameStore } from '../store/gameStore';
import { useState } from 'react';
import Swal from 'sweetalert2';

function TransactionDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { transactions } = useGameStore();
  // Find transaction by id
  const transaction = transactions.find(tx => tx.id === id);

  if (!transaction) {
    return (
      <div className="p-4">
        <div className="text-center py-12">
          <h2 className="text-xl font-semibold text-text-primary mb-2">Transaksi Tidak Ditemukan</h2>
          <p className="text-text-secondary mb-4">ID transaksi tidak valid atau sudah dihapus</p>
          <button
            onClick={() => navigate('/transactions')}
            className="px-6 py-2 bg-primary-light text-white rounded-lg hover:bg-primary-medium transition-colors"
          >
            Kembali ke Transaksi
          </button>
        </div>
      </div>
    );
  }

  const statusConfig = {
    success: {
      color: 'bg-success',
      text: 'Pembayaran Berhasil',
      description: 'Transaksi telah berhasil diproses'
    },
    pending: {
      color: 'bg-warning',
      text: 'Menunggu Pembayaran',
      description: 'Silakan lakukan pembayaran sesuai petunjuk di bawah'
    },
    failed: {
      color: 'bg-danger',
      text: 'Pembayaran Gagal',
      description: 'Transaksi gagal diproses, silakan hubungi customer service'
    }
  };

  const paymentMethods = [
    { id: 'gopay', name: 'GoPay', account: '081234567890' },
    { id: 'ovo', name: 'OVO', account: '081234567890' },
    { id: 'dana', name: 'DANA', account: '081234567890' },
    { id: 'bca', name: 'BCA', account: '1234567890' },
    { id: 'mandiri', name: 'Mandiri', account: '9876543210' }
  ];

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('id-ID', {
      weekday: 'long',
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const copyToClipboard = (text, label) => {
    navigator.clipboard.writeText(text).then(() => {
      Swal.fire({
        icon: 'success',
        title: 'Tersalin!',
        text: `${label} berhasil disalin ke clipboard`,
        timer: 1500,
        showConfirmButton: false,
      });
    });
  };

  const handlePaymentConfirmation = () => {
    Swal.fire({
      title: 'Konfirmasi Pembayaran',
      text: 'Apakah Anda sudah melakukan pembayaran?',
      icon: 'question',
      showCancelButton: true,
      confirmButtonColor: '#78A083',
      cancelButtonColor: '#6B7280',
      confirmButtonText: 'Ya, Sudah Bayar',
      cancelButtonText: 'Belum'
    }).then((result) => {
      if (result.isConfirmed) {
        Swal.fire({
          icon: 'success',
          title: 'Terima Kasih!',
          text: 'Pembayaran sedang diverifikasi, mohon tunggu konfirmasi',
          timer: 3000,
          showConfirmButton: false,
        });
      }
    });
  };

  return (
    <div className="p-4 space-y-6">
      {/* Header */}
      <header className="flex items-center space-x-4 py-6 page-background">
        <button
          onClick={() => navigate('/transactions')}
          className="p-2 hover:bg-gray-100 rounded-lg transition-colors"
          aria-label="Kembali"
        >
          <i className="fas fa-arrow-left text-xl text-text-primary"></i>
        </button>
        <div>
          <h1 className="text-xl font-bold text-primary-dark">Detail Transaksi</h1>
          <p className="text-sm text-text-secondary">ID: {transaction.id}</p>
        </div>
      </header>

      {/* Status Card */}
      <div className="bg-white shadow-md rounded-lg p-4">
        <div className="flex items-center space-x-3 mb-3">
          <div className={`w-3 h-3 rounded-full ${statusConfig[transaction.status].color}`}></div>
          <h2 className="font-semibold text-lg">{statusConfig[transaction.status].text}</h2>
        </div>
        <p className="text-text-secondary text-sm">{statusConfig[transaction.status].description}</p>
      </div>

      {/* Transaction Details */}
      <div className="bg-white shadow-md rounded-lg p-4 space-y-4">
        <h3 className="font-semibold text-lg text-primary-dark">Informasi Transaksi</h3>
        
        <div className="grid grid-cols-1 gap-3">
          <div className="flex justify-between">
            <span className="text-text-secondary">Game:</span>
            <span className="font-medium">{transaction.game}</span>
          </div>
          {transaction.productName && (
            <div className="flex justify-between">
              <span className="text-text-secondary">Produk:</span>
              <span className="font-medium">{transaction.productName}</span>
            </div>
          )}
          {transaction.gameAccount && (
            <div className="flex justify-between">
              <span className="text-text-secondary">Akun Game:</span>
              <span className="font-medium">{transaction.gameAccount}</span>
            </div>
          )}
          <div className="flex justify-between">
            <span className="text-text-secondary">Tanggal:</span>
            <span className="font-medium">{formatDate(transaction.date)}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-text-secondary">Metode Pembayaran:</span>
            <span className="font-medium">{transaction.paymentMethodName || 'Belum dipilih'}</span>
          </div>
          <div className="flex justify-between">
            <span className="text-text-secondary">Status:</span>
            <span className={`px-2 py-1 rounded text-xs font-medium text-white ${statusConfig[transaction.status].color}`}>
              {statusConfig[transaction.status].text}
            </span>
          </div>
          
          {/* Price Breakdown */}
          {(transaction.originalAmount || transaction.fee || transaction.discount) && (
            <div className="border-t pt-3 space-y-2">
              {transaction.originalAmount && (
                <div className="flex justify-between text-sm">
                  <span className="text-text-secondary">Harga Produk:</span>
                  <span>Rp {parseInt(transaction.originalAmount).toLocaleString('id-ID')}</span>
                </div>
              )}
              {transaction.fee > 0 && (
                <div className="flex justify-between text-sm">
                  <span className="text-text-secondary">Biaya Admin:</span>
                  <span>Rp {transaction.fee.toLocaleString('id-ID')}</span>
                </div>
              )}
              {transaction.discount > 0 && (
                <div className="flex justify-between text-sm text-success">
                  <span>Diskon {transaction.voucherCode ? `(${transaction.voucherCode})` : ''}:</span>
                  <span>-Rp {transaction.discount.toLocaleString('id-ID')}</span>
                </div>
              )}
            </div>
          )}
          
          <div className="border-t pt-3 flex justify-between items-center">
            <span className="text-lg font-semibold">Total Pembayaran:</span>
            <span className="text-xl font-bold text-primary-medium">
              Rp {parseInt(transaction.amount).toLocaleString('id-ID')}
            </span>
          </div>
        </div>
      </div>

      {/* Payment Instructions (only for pending status) */}
      {transaction.status === 'pending' && transaction.paymentMethod && (
        <div className="bg-white shadow-md rounded-lg p-4 space-y-4">
          <h3 className="font-semibold text-lg text-primary-dark">Petunjuk Pembayaran</h3>
          
          <div className="bg-gray-50 rounded-lg p-4 space-y-3">
            <h4 className="font-medium text-primary-dark">
              Transfer ke {transaction.paymentMethodName}
            </h4>
            
            <div className="space-y-2">
              <div className="flex justify-between items-center">
                <span className="text-text-secondary">Nomor Rekening/HP:</span>
                <div className="flex items-center space-x-2">
                  <span className="font-mono font-medium">
                    {paymentMethods.find(m => m.id === transaction.paymentMethod)?.account}
                  </span>
                  <button
                    onClick={() => copyToClipboard(
                      paymentMethods.find(m => m.id === transaction.paymentMethod)?.account,
                      'Nomor rekening'
                    )}
                    className="p-1 hover:bg-gray-200 rounded"
                  >
                    <i className="fas fa-copy text-text-secondary"></i>
                  </button>
                </div>
              </div>
              
              <div className="flex justify-between items-center">
                <span className="text-text-secondary">Jumlah Transfer:</span>
                <div className="flex items-center space-x-2">
                  <span className="font-mono font-bold text-primary-medium">
                    Rp {parseInt(transaction.amount).toLocaleString('id-ID')}
                  </span>
                  <button
                    onClick={() => copyToClipboard(transaction.amount, 'Jumlah transfer')}
                    className="p-1 hover:bg-gray-200 rounded"
                  >
                    <i className="fas fa-copy text-text-secondary"></i>
                  </button>
                </div>
              </div>
              
              {transaction.whatsapp && (
                <div className="flex justify-between items-center">
                  <span className="text-text-secondary">WhatsApp:</span>
                  <span className="font-medium">{transaction.whatsapp}</span>
                </div>
              )}
            </div>

            <div className="text-xs text-text-secondary bg-yellow-50 p-3 rounded border-l-4 border-yellow-400">
              <p className="font-medium mb-1">⚠️ Penting:</p>
              <ul className="space-y-1">
                <li>• Transfer sesuai nominal yang tertera</li>
                <li>• Pembayaran akan diverifikasi dalam 1x24 jam</li>
                <li>• Simpan bukti transfer untuk konfirmasi</li>
                <li>• Notifikasi akan dikirim ke WhatsApp yang terdaftar</li>
              </ul>
            </div>

            <button
              onClick={handlePaymentConfirmation}
              className="w-full py-3 bg-success text-white rounded-lg hover:bg-green-600 transition-colors focus:ring-2 focus:ring-green-500 focus:outline-none font-medium"
            >
              Konfirmasi Pembayaran
            </button>
          </div>
        </div>
      )}

      {/* Success Message */}
      {transaction.status === 'success' && (
        <div className="bg-green-50 border border-green-200 rounded-lg p-4">
          <div className="flex items-center space-x-2">
            <div className="w-5 h-5 bg-success rounded-full flex items-center justify-center">
              <i className="fas fa-check text-white text-xs"></i>
            </div>
            <p className="text-green-800 font-medium">
              Transaksi berhasil! Item telah dikirim ke akun game Anda.
            </p>
          </div>
        </div>
      )}

      {/* Failed Message */}
      {transaction.status === 'failed' && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <div className="text-center">
            <p className="text-red-800 font-medium mb-3">
              Transaksi gagal diproses. Silakan hubungi customer service untuk bantuan.
            </p>
            <button className="px-6 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors">
              Hubungi CS
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

export default TransactionDetail;