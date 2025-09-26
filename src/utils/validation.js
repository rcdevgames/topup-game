import * as yup from 'yup';

export const profileSchema = yup.object({
  name: yup.string().required('Nama wajib diisi'),
  phone: yup.string().matches(/^\d{10,12}$/, 'Nomor HP tidak valid (10-12 digit)').required('Nomor HP wajib diisi'),
  password: yup.string().min(6, 'Password minimal 6 karakter').nullable(),
});

export const gameAccountSchema = yup.object({
  game: yup.string().required('Game wajib dipilih'),
  gameId: yup.string().required('ID Game wajib diisi'),
  server: yup.string().required('Server wajib diisi'),
});

export const loginSchema = yup.object({
  phone: yup.string().required('Nomor HP wajib diisi'),
  password: yup.string().required('Password wajib diisi'),
});