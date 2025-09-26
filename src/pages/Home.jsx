import { useState } from 'react';
import { useGameStore } from '../store/gameStore';
import ProductCard from '../components/ProductCard';

function Home() {
  const { products, categories, loadMoreProducts, searchProducts } = useGameStore();
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('');

  const handleSearch = (e) => {
    const term = e.target.value;
    setSearchTerm(term);
    searchProducts(term);
    if (term) {
      setSelectedCategory(''); // Clear category filter when searching
    }
  };

  const handleCategoryFilter = (categoryName) => {
    setSelectedCategory(categoryName);
    setSearchTerm(''); // Clear search when filtering by category
    searchProducts(categoryName);
  };

  const clearFilters = () => {
    setSearchTerm('');
    setSelectedCategory('');
    searchProducts(''); // This will show all products
  };

  return (
    <div className="p-4 space-y-6">
      <header className="text-center py-6 page-background mb-6">
        {/* Logo Waw Store */}
        <div className="mb-4">
          <img 
            src="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjgwIiB2aWV3Qm94PSIwIDAgMjAwIDgwIiBmaWxsPSJub25lIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8ZGVmcz4KPHN0eWxlPgouY2xzLTEgeyBmaWxsOiAjMzUzNzRCOyBmb250LWZhbWlseTogSW50ZXIsIHNhbnMtc2VyaWY7IGZvbnQtc2l6ZTogMjhweDsgZm9udC13ZWlnaHQ6IDcwMDsgfQouY2xzLTIgeyBmaWxsOiAjNTA3MjdCOyBmb250LWZhbWlseTogSW50ZXIsIHNhbnMtc2VyaWY7IGZvbnQtc2l6ZTogMTJweDsgZm9udC13ZWlnaHQ6IDQwMDsgfQouaWNvbiB7IGZpbGw6ICM3OEEwODM7IH0KPC9zdHlsZT4KPC9kZWZzPgo8IS0tIEdhbWUgSWNvbiAtLT4KPHJlY3QgeD0iMTAiIHk9IjE1IiB3aWR0aD0iMzAiIGhlaWdodD0iMjAiIHJ4PSIzIiBjbGFzcz0iaWNvbiIvPgo8Y2lyY2xlIGN4PSIxOCIgY3k9IjI4IiByPSIyIiBmaWxsPSJ3aGl0ZSIvPgo8Y2lyY2xlIGN4PSIzMiIgY3k9IjI4IiByPSIyIiBmaWxsPSJ3aGl0ZSIvPgo8cmVjdCB4PSIyMCIgeT0iMzgiIHdpZHRoPSI4IiBoZWlnaHQ9IjEwIiByeD0iMiIgY2xhc3M9Imljb24iLz4KPCEtLSBUZXh0IC0tPgo8dGV4dCB4PSI1NSIgeT0iNDAiIGNsYXNzPSJjbHMtMSI+V2F3IFN0b3JlPC90ZXh0Pgo8dGV4dCB4PSI1NSIgeT0iNTYiIGNsYXNzPSJjbHMtMiI+VG9wdXAgR2FtZSBPbmxpbmU8L3RleHQ+Cjwvc3ZnPgo="
            alt="Waw Store Logo"
            className="mx-auto h-20"
          />
        </div>
        <p className="text-text-secondary">Beli diamond dan item game favorit kamu</p>
      </header>
      
      {/* Search */}
      <div className="relative">
        <i className="fas fa-search absolute left-3 top-1/2 transform -translate-y-1/2 text-text-secondary"></i>
        <input
          type="text"
          value={searchTerm}
          onChange={handleSearch}
          placeholder="Cari game atau produk..."
          className="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
          aria-label="Pencarian produk"
        />
      </div>
      
      {/* Categories */}
      <section>
        <div className="flex justify-between items-center mb-3">
          <h2 className="text-xl font-semibold text-primary-dark">Kategori</h2>
          {(selectedCategory || searchTerm) && (
            <button
              onClick={clearFilters}
              className="text-sm text-primary-light hover:text-primary-medium underline"
            >
              Lihat Semua
            </button>
          )}
        </div>
        <div className="flex overflow-x-auto space-x-3 pb-2">
          {categories.map((cat) => (
            <button
              key={cat.id}
              className={`px-4 py-2 rounded-lg whitespace-nowrap transition-colors focus:ring-2 focus:ring-blue-500 focus:outline-none ${
                selectedCategory === cat.name
                  ? 'bg-primary-medium text-white'
                  : 'bg-primary-light text-white hover:bg-primary-medium'
              }`}
              onClick={() => handleCategoryFilter(cat.name)}
            >
              {cat.name}
            </button>
          ))}
        </div>
      </section>
      
      {/* Products Grid */}
      <section>
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-xl font-semibold text-primary-dark">
            Produk {selectedCategory && `- ${selectedCategory}`}
          </h2>
          <span className="text-sm text-text-secondary">
            {products.length} produk
          </span>
        </div>
        
        {products.length === 0 ? (
          <div className="text-center py-8">
            <p className="text-text-secondary">Tidak ada produk ditemukan</p>
          </div>
        ) : (
          <>
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
              {products.map((product) => (
                <ProductCard key={product.id} product={product} />
              ))}
            </div>
            <button
              onClick={loadMoreProducts}
              className="w-full mt-6 mb-6 py-3 bg-success text-white rounded-lg hover:bg-green-600 transition-colors focus:ring-2 focus:ring-green-500 focus:outline-none font-medium"
            >
              Muat Lebih Banyak
            </button>
          </>
        )}
      </section>
    </div>
  );
}

export default Home;