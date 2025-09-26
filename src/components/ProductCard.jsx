import { useNavigate } from 'react-router-dom';

function ProductCard({ product }) {
  const navigate = useNavigate();

  const handleBuy = () => {
    navigate(`/checkout/${product.id}`);
  };

  return (
    <div className="bg-white shadow-md rounded-lg p-4 hover:shadow-lg transition-shadow flex flex-col h-full">
      <img 
        src={product.image} 
        alt={product.name} 
        className="w-full h-32 object-cover rounded mb-2"
        onError={(e) => {
          e.target.src = 'https://images.unsplash.com/photo-1511512578047-dfb367046420?w=150&h=100&fit=crop&crop=center';
        }}
      />
      <div className="flex-1 flex flex-col">
        <h3 className="font-semibold text-sm mb-1">{product.name}</h3>
        <p className="text-xs text-text-secondary mb-2 line-clamp-2 flex-1">{product.description}</p>
        <p className="text-primary-medium font-bold mb-3">Rp {parseInt(product.price).toLocaleString('id-ID')}</p>
      </div>
      <button
        onClick={handleBuy}
        className="w-full py-2 bg-primary-light text-white rounded hover:bg-primary-medium transition-colors focus:ring-2 focus:ring-blue-500 focus:outline-none mt-auto"
        aria-label={`Beli ${product.name}`}
      >
        Beli
      </button>
    </div>
  );
}

export default ProductCard;