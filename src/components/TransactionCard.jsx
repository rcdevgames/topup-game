import { useNavigate } from 'react-router-dom';

function TransactionCard({ transaction }) {
  const navigate = useNavigate();
  
  const statusColors = {
    success: 'bg-success text-white',
    pending: 'bg-warning text-white',
    failed: 'bg-danger text-white',
  };

  const statusText = {
    success: 'Berhasil',
    pending: 'Pending',
    failed: 'Gagal',
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('id-ID', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  const handleCardClick = () => {
    navigate(`/transactions/${transaction.id}`);
  };

  return (
    <div 
      className="bg-white shadow-md rounded-lg p-4 hover:shadow-lg transition-shadow cursor-pointer"
      onClick={handleCardClick}
      role="button"
      tabIndex={0}
      onKeyDown={(e) => {
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault();
          handleCardClick();
        }
      }}
      aria-label={`Lihat detail transaksi ${transaction.id}`}>
      <div className="flex justify-between items-center mb-2">
        <h3 className="font-semibold text-lg">{transaction.game}</h3>
        <span className={`px-3 py-1 rounded-full text-xs font-medium ${statusColors[transaction.status]}`}>
          {statusText[transaction.status]}
        </span>
      </div>
      <div className="space-y-1">
        <p className="text-sm text-text-secondary">
          <span className="font-medium">ID Transaksi:</span> {transaction.id}
        </p>
        <p className="text-sm text-text-secondary">
          <span className="font-medium">Tanggal:</span> {formatDate(transaction.date)}
        </p>
        <p className="text-lg font-bold text-primary-medium mt-2">
          Rp {parseInt(transaction.amount).toLocaleString('id-ID')}
        </p>
      </div>
      
      {/* Click indicator */}
      <div className="flex items-center justify-end mt-2">
        <i className="fas fa-chevron-right text-text-secondary"></i>
      </div>
    </div>
  );
}

export default TransactionCard;