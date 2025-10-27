'use client';

import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import AdminLayout from '@/components/AdminLayout';
import { representativeAPI } from '@/lib/api';
import type { UpdateRepresentativeData, RepresentativeConfig } from '@/types';

export default function SettingsPage() {
  const [config, setConfig] = useState<RepresentativeConfig | null>(null);
  const [loading, setLoading] = useState(false);
  const [loadingData, setLoadingData] = useState(true);
  const { register, handleSubmit, formState: { errors }, reset } = useForm<UpdateRepresentativeData>();

  useEffect(() => {
    loadConfig();
  }, []);

  const loadConfig = async () => {
    setLoadingData(true);
    try {
      const data = await representativeAPI.getConfig();
      setConfig(data);
      reset({
        name: data.name,
        position: data.position,
      });
    } catch (error) {
      console.error('Failed to load representative config:', error);
      alert('Gagal memuat konfigurasi perwakilan');
    } finally {
      setLoadingData(false);
    }
  };

  const onSubmit = async (data: UpdateRepresentativeData) => {
    setLoading(true);
    try {
      const updatedConfig = await representativeAPI.updateConfig(data);
      setConfig(updatedConfig);
      alert('Konfigurasi perwakilan berhasil diperbarui!');
    } catch (error: any) {
      console.error('Failed to update representative config:', error);
      alert(error.response?.data?.error || 'Gagal memperbarui konfigurasi perwakilan');
    } finally {
      setLoading(false);
    }
  };

  return (
    <AdminLayout>
      <div className="px-4 py-6 sm:px-0">
        <h2 className="text-2xl font-bold text-gray-800 mb-6">Pengaturan Sistem</h2>

        <div className="max-w-2xl">
          {/* Representative Config Card */}
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4">
              Konfigurasi Perwakilan Penandatangan
            </h3>

            <p className="text-sm text-gray-600 mb-6">
              Pengaturan nama dan jabatan perwakilan yang akan ditampilkan pada dokumen perjalanan dinas (Nota Permintaan dan Berita Acara).
            </p>

            {loadingData ? (
              <div className="flex justify-center py-8">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
              </div>
            ) : (
              <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                {/* Name Field */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Nama Perwakilan <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    {...register('name', { required: 'Nama perwakilan wajib diisi' })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="Contoh: M. MACHFUD HIDAYAT"
                  />
                  {errors.name && (
                    <p className="mt-1 text-sm text-red-600">{errors.name.message}</p>
                  )}
                </div>

                {/* Position Field */}
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Jabatan <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    {...register('position', { required: 'Jabatan wajib diisi' })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="Contoh: Vice President"
                  />
                  {errors.position && (
                    <p className="mt-1 text-sm text-red-600">{errors.position.message}</p>
                  )}
                </div>

                {/* Current Values Info */}
                {config && (
                  <div className="bg-blue-50 border border-blue-200 rounded-md p-4 mt-6">
                    <h4 className="text-sm font-semibold text-blue-900 mb-2">
                      Nilai Saat Ini:
                    </h4>
                    <dl className="space-y-1">
                      <div className="flex">
                        <dt className="text-sm text-blue-700 w-24">Nama:</dt>
                        <dd className="text-sm text-blue-900 font-medium">{config.name}</dd>
                      </div>
                      <div className="flex">
                        <dt className="text-sm text-blue-700 w-24">Jabatan:</dt>
                        <dd className="text-sm text-blue-900 font-medium">{config.position}</dd>
                      </div>
                    </dl>
                  </div>
                )}

                {/* Submit Button */}
                <div className="flex justify-end pt-4">
                  <button
                    type="submit"
                    disabled={loading}
                    className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    {loading ? 'Menyimpan...' : 'Simpan Perubahan'}
                  </button>
                </div>
              </form>
            )}
          </div>

          {/* Info Card */}
          <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mt-6">
            <div className="flex">
              <div className="flex-shrink-0">
                <svg className="h-5 w-5 text-yellow-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                </svg>
              </div>
              <div className="ml-3">
                <h3 className="text-sm font-medium text-yellow-800">Perhatian</h3>
                <div className="mt-2 text-sm text-yellow-700">
                  <p>
                    Perubahan konfigurasi ini akan berlaku untuk semua permohonan perjalanan dinas baru yang dibuat setelah perubahan disimpan. Dokumen yang sudah dibuat sebelumnya tidak akan terpengaruh.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </AdminLayout>
  );
}
