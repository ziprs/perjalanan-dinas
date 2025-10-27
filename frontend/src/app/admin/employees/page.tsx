'use client';

import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import AdminLayout from '@/components/AdminLayout';
import { employeeAPI, positionAPI } from '@/lib/api';
import type { Employee, Position, CreateEmployeeData } from '@/types';

export default function EmployeesPage() {
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [positions, setPositions] = useState<Position[]>([]);
  const [loading, setLoading] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const { register, handleSubmit, formState: { errors }, reset } = useForm<CreateEmployeeData>();

  useEffect(() => {
    loadEmployees();
    loadPositions();
  }, []);

  const loadEmployees = async () => {
    try {
      const data = await employeeAPI.getAll();
      setEmployees(data);
    } catch (error) {
      console.error('Failed to load employees:', error);
      alert('Gagal memuat data karyawan');
    }
  };

  const loadPositions = async () => {
    try {
      const data = await positionAPI.getAll();
      setPositions(data);
    } catch (error) {
      console.error('Failed to load positions:', error);
      alert('Gagal memuat data jabatan');
    }
  };

  const onSubmit = async (data: CreateEmployeeData) => {
    setLoading(true);
    try {
      if (editingId) {
        await employeeAPI.update(editingId, data);
        alert('Karyawan berhasil diupdate!');
      } else {
        await employeeAPI.create(data);
        alert('Karyawan berhasil ditambahkan!');
      }
      reset();
      setEditingId(null);
      loadEmployees();
    } catch (error: any) {
      console.error('Failed to save employee:', error);
      alert(error.response?.data?.error || 'Gagal menyimpan data karyawan');
    } finally {
      setLoading(false);
    }
  };

  const handleEdit = (employee: Employee) => {
    setEditingId(employee.id);
    reset({
      nip: employee.nip,
      name: employee.name,
      position_id: employee.position_id,
    });
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Apakah Anda yakin ingin menghapus karyawan ini?')) return;

    try {
      await employeeAPI.delete(id);
      alert('Karyawan berhasil dihapus!');
      loadEmployees();
    } catch (error) {
      console.error('Failed to delete employee:', error);
      alert('Gagal menghapus karyawan');
    }
  };

  const handleCancel = () => {
    setEditingId(null);
    reset();
  };

  return (
    <AdminLayout>
      <div className="px-4 py-6 sm:px-0">
        <h2 className="text-2xl font-bold text-gray-800 mb-6">Kelola Karyawan</h2>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Form */}
          <div className="lg:col-span-1">
            <div className="bg-white rounded-lg shadow p-6">
              <h3 className="text-lg font-semibold text-gray-800 mb-4">
                {editingId ? 'Edit Karyawan' : 'Tambah Karyawan Baru'}
              </h3>
              <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    NIP <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    {...register('nip', { required: 'NIP harus diisi' })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                    placeholder="Contoh: 123456789"
                  />
                  {errors.nip && (
                    <p className="mt-1 text-sm text-red-600">{errors.nip.message}</p>
                  )}
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Nama <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="text"
                    {...register('name', { required: 'Nama harus diisi' })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                    placeholder="Contoh: John Doe"
                  />
                  {errors.name && (
                    <p className="mt-1 text-sm text-red-600">{errors.name.message}</p>
                  )}
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Jabatan <span className="text-red-500">*</span>
                  </label>
                  <select
                    {...register('position_id', {
                      required: 'Jabatan harus dipilih',
                      valueAsNumber: true
                    })}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="">-- Pilih Jabatan --</option>
                    {positions.map((position) => (
                      <option key={position.id} value={position.id}>
                        {position.title} ({position.code})
                      </option>
                    ))}
                  </select>
                  {errors.position_id && (
                    <p className="mt-1 text-sm text-red-600">{errors.position_id.message}</p>
                  )}
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
            </div>
          </div>

          {/* List */}
          <div className="lg:col-span-2">
            <div className="bg-white rounded-lg shadow overflow-hidden">
              <div className="px-6 py-4 border-b border-gray-200">
                <h3 className="text-lg font-semibold text-gray-800">Daftar Karyawan</h3>
              </div>
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">NIP</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Nama</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Jabatan</th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Aksi</th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {employees.map((employee) => (
                      <tr key={employee.id}>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{employee.nip}</td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{employee.name}</td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                          {employee.position ? employee.position.title : '-'}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm">
                          <button
                            onClick={() => handleEdit(employee)}
                            className="text-blue-600 hover:text-blue-900 mr-3"
                          >
                            Edit
                          </button>
                          <button
                            onClick={() => handleDelete(employee.id)}
                            className="text-red-600 hover:text-red-900"
                          >
                            Hapus
                          </button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
                {employees.length === 0 && (
                  <div className="text-center py-8 text-gray-500">
                    Belum ada data karyawan
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
