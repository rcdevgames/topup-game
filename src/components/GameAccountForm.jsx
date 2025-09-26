import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import { gameAccountSchema } from '../utils/validation';

function GameAccountForm({ onSubmit, initialData = {}, isEdit = false }) {
  const { register, handleSubmit, formState: { errors }, reset } = useForm({
    resolver: yupResolver(gameAccountSchema),
    defaultValues: initialData,
  });

  const handleFormSubmit = (data) => {
    onSubmit(data);
    if (!isEdit) {
      reset(); // Clear form after successful add
    }
  };

  const gameOptions = [
    'Mobile Legends',
    'Free Fire',
    'PUBG Mobile',
    'Genshin Impact',
    'Arena of Valor',
    'Call of Duty Mobile'
  ];

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="bg-white/80 backdrop-blur-sm border border-white/20 shadow-lg rounded-lg p-4 space-y-4">
      <h3 className="text-lg font-semibold text-primary-dark mb-4">
        {isEdit ? 'Edit Akun Game' : 'Tambah Akun Game'}
      </h3>
      
      <div>
        <label className="block text-sm font-medium text-text-primary mb-2">Game</label>
        <select 
          {...register('game')} 
          className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
        >
          <option value="">Pilih Game</option>
          {gameOptions.map((game) => (
            <option key={game} value={game}>{game}</option>
          ))}
        </select>
        {errors.game && <p className="text-danger text-xs mt-1">{errors.game.message}</p>}
      </div>
      
      <div>
        <label className="block text-sm font-medium text-text-primary mb-2">ID Game</label>
        <input 
          {...register('gameId')} 
          className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
          placeholder="Masukkan ID Game"
        />
        {errors.gameId && <p className="text-danger text-xs mt-1">{errors.gameId.message}</p>}
      </div>
      
      <div>
        <label className="block text-sm font-medium text-text-primary mb-2">Server</label>
        <input 
          {...register('server')} 
          className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 focus:outline-none"
          placeholder="Contoh: 2001, Asia, Global"
        />
        {errors.server && <p className="text-danger text-xs mt-1">{errors.server.message}</p>}
      </div>
      
      <button
        type="submit"
        className="w-full py-3 bg-primary-light text-white rounded-lg hover:bg-primary-medium transition-colors focus:ring-2 focus:ring-blue-500 focus:outline-none font-medium"
      >
        {isEdit ? 'Update Akun' : 'Tambah Akun'}
      </button>
    </form>
  );
}

export default GameAccountForm;