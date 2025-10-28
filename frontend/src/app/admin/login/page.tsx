'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { authAPI } from '@/lib/api';
import type { LoginData } from '@/types';

export default function AdminLogin() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const { register, handleSubmit, formState: { errors } } = useForm<LoginData>();

  const onSubmit = async (data: LoginData) => {
    setLoading(true);
    try {
      const response = await authAPI.login(data);
      localStorage.setItem('token', response.token);
      localStorage.setItem('username', response.username);
      alert('Login berhasil!');
      router.push('/admin/dashboard');
    } catch (error: any) {
      console.error('Login failed:', error);
      alert(error.response?.data?.error || 'Login gagal');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-800 to-gray-900 flex items-center justify-center py-12 px-4">
      <div className="max-w-md w-full">
        <div className="bg-white rounded-lg shadow-2xl p-8">
          {/* Logo Section */}
          <div className="flex justify-center mb-6">
            <img
              src="/logo-digibank.png"
              alt="Divisi Digital Banking"
              className="h-24 w-auto object-contain"
            />
          </div>

          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-gray-800 mb-2">Admin Login</h1>
            <p className="text-gray-600">Sistem Perjalanan Dinas</p>
          </div>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Username
              </label>
              <input
                type="text"
                {...register('username', { required: 'Username harus diisi' })}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Masukkan username"
              />
              {errors.username && (
                <p className="mt-1 text-sm text-red-600">{errors.username.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Password
              </label>
              <input
                type="password"
                {...register('password', { required: 'Password harus diisi' })}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Masukkan password"
              />
              {errors.password && (
                <p className="mt-1 text-sm text-red-600">{errors.password.message}</p>
              )}
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed font-medium transition-colors"
            >
              {loading ? 'Memproses...' : 'Login'}
            </button>

            <div className="text-center">
              <a href="/spd" className="text-sm text-blue-600 hover:text-blue-700">
                Kembali ke Halaman Utama
              </a>
            </div>
          </form>
        </div>

        <div className="mt-6 text-center text-gray-400 text-sm">
          <p>@UncleJSS</p>
        </div>
      </div>
    </div>
  );
}
