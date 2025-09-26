import { useGameStore } from '../store/gameStore';
import TransactionCard from '../components/TransactionCard';

function Transactions() {
  const { transactions, loadMoreTransactions } = useGameStore();

  return (
    <div className="p-4 space-y-6">
      <header className="py-6 page-background">
        <h1 className="text-2xl font-bold text-primary-dark mb-2">Riwayat Transaksi</h1>
        <p className="text-text-secondary">Lihat semua transaksi topup game kamu</p>
      </header>
      
      {transactions.length === 0 ? (
        <div className="text-center py-12">
          <div className="mb-4">
            <i className="fas fa-file-invoice text-6xl text-text-secondary"></i>
          </div>
          <h3 className="text-lg font-medium text-text-primary mb-2">Belum ada transaksi</h3>
          <p className="text-text-secondary">Transaksi topup game kamu akan muncul di sini</p>
        </div>
      ) : (
        <>
          <div className="space-y-4">
            {transactions.map((tx) => (
              <TransactionCard key={tx.id} transaction={tx} />
            ))}
          </div>
          
          <button
            onClick={loadMoreTransactions}
            className="w-full py-3 bg-success text-white rounded-lg hover:bg-green-600 transition-colors focus:ring-2 focus:ring-green-500 focus:outline-none font-medium"
          >
            Muat Lebih Banyak
          </button>
        </>
      )}
    </div>
  );
}

export default Transactions;