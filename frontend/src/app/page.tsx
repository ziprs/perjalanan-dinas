'use client';

import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { employeeAPI, positionCodeAPI, travelRequestAPI, pdfAPI, cityAPI } from '@/lib/api';
import type { Employee, CreateTravelRequestData, City } from '@/types';

export default function Home() {
  const [employees, setEmployees] = useState<Employee[]>([]);
  const [cities, setCities] = useState<City[]>([]);
  const [selectedEmployeeIds, setSelectedEmployeeIds] = useState<number[]>([]);
  const [nipInput, setNipInput] = useState('');
  const [loading, setLoading] = useState(false);
  const [submittedRequestId, setSubmittedRequestId] = useState<number | null>(null);

  const { register, handleSubmit, formState: { errors }, reset, watch } = useForm<Omit<CreateTravelRequestData, 'employee_ids'>>({
    defaultValues: {
      departure_place: 'Surabaya'
    }
  });

  const departureDate = watch('departure_date');
  const returnDate = watch('return_date');
  const destinationType = watch('destination_type');

  useEffect(() => {
    loadEmployees();
    loadCities();
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

  const loadCities = async () => {
    try {
      const data = await cityAPI.getAll();
      setCities(data);
    } catch (error) {
      console.error('Failed to load cities:', error);
      alert('Gagal memuat data kota');
    }
  };

  const calculateDuration = () => {
    if (departureDate && returnDate) {
      const start = new Date(departureDate);
      const end = new Date(returnDate);
      const diffTime = Math.abs(end.getTime() - start.getTime());
      const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
      return diffDays + 1;
    }
    return 0;
  };

  const handleAddEmployeeByNIP = () => {
    const trimmedNip = nipInput.trim();
    if (!trimmedNip) {
      return;
    }

    const employee = employees.find(emp => emp.nip === trimmedNip);
    if (!employee) {
      alert(`Karyawan dengan NIP "${trimmedNip}" tidak ditemukan`);
      return;
    }

    if (selectedEmployeeIds.includes(employee.id)) {
      alert(`Karyawan ${employee.name} sudah ditambahkan`);
      setNipInput('');
      return;
    }

    setSelectedEmployeeIds(prev => [...prev, employee.id]);
    setNipInput('');
  };

  const handleRemoveEmployee = (employeeId: number) => {
    setSelectedEmployeeIds(prev => prev.filter(id => id !== employeeId));
  };

  const handleNipKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      handleAddEmployeeByNIP();
    }
  };

  const onSubmit = async (data: Omit<CreateTravelRequestData, 'employee_ids'>) => {
    if (selectedEmployeeIds.length === 0) {
      alert('Pilih minimal 1 karyawan');
      return;
    }

    setLoading(true);
    try {
      const response = await travelRequestAPI.create({
        ...data,
        employee_ids: selectedEmployeeIds,
        departure_place: data.departure_place || 'Surabaya',
      });

      alert('Perjalanan dinas berhasil dibuat!');
      setSubmittedRequestId(response.id);
      setSelectedEmployeeIds([]);
      reset({
        departure_place: 'Surabaya'
      });
    } catch (error: any) {
      console.error('Failed to create travel request:', error);
      alert(error.response?.data?.error || 'Gagal membuat perjalanan dinas');
    } finally {
      setLoading(false);
    }
  };

  const selectedEmployees = employees.filter(emp => selectedEmployeeIds.includes(emp.id));

  // Filter cities based on selected destination_type
  const filteredCities = destinationType
    ? cities.filter(city => city.destination_type === destinationType)
    : [];

  return (
    <main className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 py-12 px-4">
      <div className="max-w-4xl mx-auto">
        <div className="bg-white rounded-lg shadow-xl p-8">
          {/* Logo Section */}
          <div className="flex justify-center mb-6">
            <img
              src="/logo-digibank.png"
              alt="Divisi Digital Banking"
              className="h-24 w-auto object-contain"
            />
          </div>

          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-gray-800 mb-2">
              INPUT SPD
            </h1>
            <p className="text-gray-600">
              Form Pengajuan SPD
            </p>
          </div>

          {submittedRequestId && (
            <div className="mb-6 p-4 bg-green-50 border border-green-200 rounded-lg">
              <h3 className="font-semibold text-green-800 mb-2">✅ Berhasil!</h3>
              <p className="text-green-700 mb-3">
                Perjalanan dinas telah dibuat. Dokumen siap diunduh:
              </p>
              <a
                href={pdfAPI.downloadCombined(submittedRequestId)}
                className="inline-block px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium transition-colors shadow-md"
                target="_blank"
                rel="noopener noreferrer"
              >
                📥 Unduh Dokumen (Nota Permintaan & Berita Acara)
              </a>
              <p className="mt-3 text-sm text-green-600">
                💡 File PDF berisi 2 halaman: Halaman 1 (Nota Permintaan) dan Halaman 2 (Berita Acara)
              </p>
            </div>
          )}

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            {/* Employee Selection by NIP */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Input NIP Karyawan <span className="text-red-500">*</span>
                <span className="ml-2 text-xs text-gray-500">(Dapat menginput lebih dari 1)</span>
              </label>
              <div className="flex gap-2">
                <input
                  type="text"
                  value={nipInput}
                  onChange={(e) => setNipInput(e.target.value)}
                  onKeyPress={handleNipKeyPress}
                  placeholder="Masukkan NIP karyawan"
                  className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  disabled={employees.length === 0}
                />
                <button
                  type="button"
                  onClick={handleAddEmployeeByNIP}
                  disabled={employees.length === 0 || !nipInput.trim()}
                  className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
                >
                  Tambah
                </button>
              </div>
              {employees.length === 0 && (
                <p className="mt-2 text-sm text-gray-500">
                  Belum ada data karyawan. Silakan tambahkan melalui Admin Panel.
                </p>
              )}
              {selectedEmployeeIds.length === 0 && employees.length > 0 && (
                <p className="mt-2 text-sm text-amber-600">
                  ⚠ Input minimal 1 NIP karyawan untuk melanjutkan
                </p>
              )}
            </div>

            {/* Selected Employees Summary */}
            {selectedEmployees.length > 0 && (
              <div className="p-4 bg-blue-50 rounded-lg border border-blue-200">
                <h3 className="font-semibold text-blue-900 mb-3">
                  Karyawan Terpilih ({selectedEmployees.length}):
                </h3>
                <div className="space-y-2">
                  {selectedEmployees.map((emp) => (
                    <div key={emp.id} className="flex justify-between items-center bg-white p-3 rounded shadow-sm">
                      <div className="flex-1">
                        <div className="font-medium text-gray-900">{emp.name}</div>
                        <div className="text-sm text-gray-600">
                          {emp.position ? emp.position.title : '-'} • NIP: {emp.nip}
                        </div>
                      </div>
                      <button
                        type="button"
                        onClick={() => handleRemoveEmployee(emp.id)}
                        className="ml-4 px-3 py-1 text-red-600 hover:bg-red-50 rounded text-sm font-medium transition-colors"
                      >
                        Hapus
                      </button>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Purpose */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Maksud Perjalanan Dinas <span className="text-red-500">*</span>
              </label>
              <textarea
                {...register('purpose', { required: 'Maksud perjalanan dinas harus diisi' })}
                rows={3}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Contoh: Melakukan kunjungan kerja ke kantor cabang..."
              />
              {errors.purpose && (
                <p className="mt-1 text-sm text-red-600">{errors.purpose.message}</p>
              )}
            </div>

            {/* Departure Place */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Tempat Berangkat
                <span className="ml-2 text-xs font-normal text-gray-500">(Dapat diubah jika berbeda)</span>
              </label>
              <input
                type="text"
                {...register('departure_place')}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Surabaya"
              />
              <p className="mt-1 text-xs text-gray-500">
                💡 Default: Surabaya (dapat diubah sesuai kebutuhan)
              </p>
            </div>

            {/* Destination Type */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Jenis Tujuan <span className="text-red-500">*</span>
              </label>
              <select
                {...register('destination_type', { required: 'Jenis tujuan harus dipilih' })}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="">-- Pilih Jenis Tujuan --</option>
                <option value="in_province">Dalam Provinsi</option>
                <option value="outside_province">Luar Provinsi</option>
                <option value="abroad">Luar Negeri</option>
              </select>
              {errors.destination_type && (
                <p className="mt-1 text-sm text-red-600">{errors.destination_type.message}</p>
              )}
            </div>

            {/* Destination */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Tempat Tujuan <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                list="city-options"
                {...register('destination', { required: 'Tempat tujuan harus diisi' })}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder={
                  destinationType
                    ? "Ketik untuk mencari atau pilih dari dropdown"
                    : "Pilih jenis tujuan terlebih dahulu"
                }
                disabled={!destinationType}
              />
              <datalist id="city-options">
                {filteredCities.map((city, index) => (
                  <option key={`${city.name}-${index}`} value={city.name} />
                ))}
              </datalist>
              {errors.destination && (
                <p className="mt-1 text-sm text-red-600">{errors.destination.message}</p>
              )}
              {destinationType && filteredCities.length === 0 && (
                <p className="mt-2 text-sm text-amber-600">
                  Belum ada kota untuk jenis tujuan ini. Silakan hubungi admin.
                </p>
              )}
              {destinationType && filteredCities.length > 0 && (
                <p className="mt-2 text-sm text-gray-500">
                  {filteredCities.length} kota tersedia untuk dipilih
                </p>
              )}
            </div>

            {/* Dates */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Tanggal Berangkat <span className="text-red-500">*</span>
                </label>
                <input
                  type="date"
                  {...register('departure_date', { required: 'Tanggal berangkat harus diisi' })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
                {errors.departure_date && (
                  <p className="mt-1 text-sm text-red-600">{errors.departure_date.message}</p>
                )}
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Tanggal Kembali <span className="text-red-500">*</span>
                </label>
                <input
                  type="date"
                  {...register('return_date', { required: 'Tanggal kembali harus diisi' })}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
                {errors.return_date && (
                  <p className="mt-1 text-sm text-red-600">{errors.return_date.message}</p>
                )}
              </div>
            </div>

            {/* Duration (auto-calculated) */}
            {departureDate && returnDate && (
              <div className="p-3 bg-blue-50 rounded-lg border border-blue-200">
                <p className="text-sm text-blue-800">
                  <span className="font-semibold">Lama Perjalanan:</span> {calculateDuration()} hari
                </p>
              </div>
            )}

            {/* Transportation */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Angkutan yang Digunakan <span className="text-red-500">*</span>
              </label>
              <select
                {...register('transportation', { required: 'Angkutan harus dipilih' })}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="">-- Pilih Angkutan --</option>
                <option value="Kendaraan Umum">Kendaraan Umum</option>
                <option value="Kendaraan Dinas">Kendaraan Dinas</option>
                <option value="Pesawat">Pesawat</option>
                <option value="Kereta Api">Kereta Api</option>
              </select>
              {errors.transportation && (
                <p className="mt-1 text-sm text-red-600">{errors.transportation.message}</p>
              )}
            </div>

            {/* Submit Button */}
            <div className="flex gap-4 pt-4">
              <button
                type="submit"
                disabled={loading}
                className="flex-1 px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed font-medium transition-colors"
              >
                {loading ? 'Memproses...' : 'Buat Perjalanan Dinas'}
              </button>

              <a
                href="/spd/admin/login"
                className="px-6 py-3 bg-gray-600 text-white rounded-lg hover:bg-gray-700 font-medium transition-colors"
              >
                Admin Login
              </a>
            </div>
          </form>
        </div>

        <div className="mt-6 text-center text-gray-600 text-sm">
          <p>UncleJSS</p>
        </div>
      </div>
    </main>
  );
}
