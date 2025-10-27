'use client';

import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import AdminLayout from '@/components/AdminLayout';
import { positionCodeAPI } from '@/lib/api';
import type { PositionCode, CreatePositionCodeData } from '@/types';

export default function PositionCodesPage() {
  const [positionCodes, setPositionCodes] = useState<PositionCode[]>([]);
  const [loading, setLoading] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const { register, handleSubmit, formState: { errors }, reset } = useForm<CreatePositionCodeData>();

  useEffect(() => {
    loadPositionCodes();
  }, []);

  const loadPositionCodes = async () => {
    try {
      const data = await positionCodeAPI.getAll();
      setPositionCodes(data);
    } catch (error) {
      console.error('Failed to load position codes:', error);
      alert('Gagal memuat data kode jabatan');
    }
  };

  const onSubmit = async (data: CreatePositionCodeData) => {
    setLoading(true);
    try {
      if (editingId) {
        await positionCodeAPI.update(editingId, data);
        alert('Kode jabatan berhasil diupdate!');
      } else {
        await positionCodeAPI.create(data);
        alert('Kode jabatan berhasil ditambahkan!');
      }
      reset();
      setEditingId(null);
      loadPositionCodes();
    } catch (error: any) {
      console.error('Failed to save position code:', error);
      alert(error.response?.data?.error || 'Gagal menyimpan kode jabatan');
    } finally {
      setLoading(false);
    }
  };

  const handleEdit = (positionCode: PositionCode) => {
    setEditingId(positionCode.id);
    reset({
      position: positionCode.position,
      code: positionCode.code,
    });
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Apakah Anda yakin ingin menghapus kode jabatan ini?')) return;

    try {
      await positionCodeAPI.delete(id);
      alert('Kode jabatan berhasil dihapus!');
      loadPositionCodes();
    } catch (error) {
      console.error('Failed to delete position code:', error);
      alert('Gagal menghapus kode jabatan');
    }
  };

  const handleCancel = () => {
    setEditingId(null);
    reset();
  };

  return (
    <AdminLayout>
      <div className="px-4 py-6 sm:px-0">
        <h2 className="text-2xl font-bold text-gray-800 mb-6">Pengkodean Jabatan</h2>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Form */}
          <div className="lg:col-span-1">
            <div className="bg-white rounded-lg shadow p-6">
              <h3 className="text-lg font-semibold text-gray-800 mb-4">
                {editingId ? 'Edit Kode Jabatan' : 'Tambah Kode Jabatan'}
              </h3>
              <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Nama Jabatan <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    {...register('position', { required: 'Nama jabatan harus diisi' })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                    placeholder="Contoh: Manager"
                  />
                  {errors.position && (
                    <p className="mt-1 text-sm text-red-600">{errors.position.message}</p>
                  )}
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Kode Jabatan <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    {...register('code', { required: 'Kode jabatan harus diisi' })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                    placeholder="Contoh: MNG"
                  />
                  {errors.code && (
                    <p className="mt-1 text-sm text-red-600">{errors.code.message}</p>
                  )}
                  <p className="mt-1 text-xs text-gray-500">
                    Kode ini akan digunakan dalam format nomor: 064/xxxx/DIB/[KODE]/NOTA
                  </p>
                </div>

                <div className="flex gap-2">
                  <button
                    type="submit"
                    disabled={loading}
                    className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400"
                  >
                    {loading ? 'Menyimpan...' : editingId ? 'Update' : 'Tambah'}
                  </button>
                  {editingId && (
                    <button
                      type="button"
                      onClick={handleCancel}
                      className="px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700"
                    >
                      Batal
                    </button>
                  )}
                </div>
              </form>

              <div className="mt-6 p-4 bg-blue-50 rounded-lg">
                <h4 className="text-sm font-semibold text-blue-900 mb-2">Informasi</h4>
                <p className="text-xs text-blue-800">
                  Pastikan setiap jabatan memiliki kode yang unik. Kode jabatan akan digunakan dalam penomoran dokumen perjalanan dinas.
                </p>
              </div>
            </div>
          </div>

          {/* List */}
          <div className="lg:col-span-2">
            <div className="bg-white rounded-lg shadow overflow-hidden">
              <div className="px-6 py-4 border-b border-gray-200">
                <h3 className="text-lg font-semibold text-gray-800">Daftar Kode Jabatan</h3>
              </div>
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Nama Jabatan</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Kode</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Format Nomor</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Aksi</th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {positionCodes.map((pc) => (
                      <tr key={pc.id}>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{pc.position}</td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm">
                          <span className="px-2 py-1 bg-purple-100 text-purple-800 rounded font-mono">
                            {pc.code}
                          </span>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600 font-mono">
                          064/xxxx/DIB/{pc.code}/NOTA
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm">
                          <button
                            onClick={() => handleEdit(pc)}
                            className="text-blue-600 hover:text-blue-900 mr-3"
                          >
                            Edit
                          </button>
                          <button
                            onClick={() => handleDelete(pc.id)}
                            className="text-red-600 hover:text-red-900"
                          >
                            Hapus
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
                {positionCodes.length === 0 && (
                  <div className="text-center py-8 text-gray-500">
                    Belum ada data kode jabatan
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </AdminLayout>
  );
}
