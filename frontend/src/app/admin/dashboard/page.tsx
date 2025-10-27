'use client';

import { useEffect, useState } from 'react';
import AdminLayout from '@/components/AdminLayout';
import { travelRequestAPI, employeeAPI } from '@/lib/api';
import Link from 'next/link';

export default function AdminDashboard() {
  const [stats, setStats] = useState({
    totalRequests: 0,
    totalEmployees: 0,
    totalPositions: 21,
  });

  useEffect(() => {
    loadStats();
  }, []);

  const loadStats = async () => {
    try {
      const [requests, employees] = await Promise.all([
        travelRequestAPI.getAll(),
        employeeAPI.getAll(),
      ]);

      setStats({
        totalRequests: requests.length,
        totalEmployees: employees.length,
        totalPositions: 21,
      });
    } catch (error) {
      console.error('Failed to load stats:', error);
    }
  };

  const statCards = [
    {
      title: 'Total Perjalanan Dinas',
      value: stats.totalRequests,
      color: 'bg-blue-500',
      icon: '📋',
    },
    {
      title: 'Total Karyawan',
      value: stats.totalEmployees,
      color: 'bg-green-500',
      icon: '👥',
    },
    {
      title: 'Total Jabatan',
      value: stats.totalPositions,
      color: 'bg-purple-500',
      icon: '💼',
    },
  ];

  return (
    <AdminLayout>
      <div className="px-4 py-6 sm:px-0">
        <h2 className="text-2xl font-bold text-gray-800 mb-6">Dashboard</h2>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
          {statCards.map((stat, index) => (
            <div
              key={index}
              className={`${stat.color} rounded-lg shadow-lg p-6 text-white`}
            >
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm opacity-90">{stat.title}</p>
                  <p className="text-3xl font-bold mt-2">{stat.value}</p>
                </div>
                <div className="text-4xl opacity-80">{stat.icon}</div>
              </div>
            </div>
          ))}
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4">Menu Utama</h3>
            <div className="space-y-3">
              <Link
                href="/admin/employees"
                className="block p-4 bg-blue-50 hover:bg-blue-100 rounded-lg transition-colors"
              >
                <div className="font-medium text-blue-900">Kelola Karyawan</div>
                <div className="text-sm text-blue-700">Tambah, edit, hapus data karyawan</div>
              </Link>
              <Link
                href="/admin/monitoring-iuran"
                className="block p-4 bg-purple-50 hover:bg-purple-100 rounded-lg transition-colors"
              >
                <div className="font-medium text-purple-900">Monitoring Iuran</div>
                <div className="text-sm text-purple-700">Lihat rekap iuran perjalanan dinas per pegawai</div>
              </Link>
              <Link
                href="/admin/travel-requests"
                className="block p-4 bg-green-50 hover:bg-green-100 rounded-lg transition-colors"
              >
                <div className="font-medium text-green-900">Daftar Perjalanan Dinas</div>
                <div className="text-sm text-green-700">Lihat dan kelola perjalanan dinas</div>
              </Link>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4">Informasi</h3>
            <div className="space-y-3 text-sm text-gray-600">
              <div className="p-3 bg-gray-50 rounded">
                <p className="font-medium text-gray-800">Sistem Jabatan:</p>
                <p className="mt-1">Sistem menggunakan 21 jabatan tetap dengan tarif iuran berbeda berdasarkan jenis perjalanan (dalam provinsi, luar provinsi, luar negeri).</p>
              </div>
              <div className="p-3 bg-gray-50 rounded">
                <p className="font-medium text-gray-800">Perhitungan Iuran:</p>
                <p className="mt-1">Total Iuran = Jumlah Hari × Tarif Jabatan × Jumlah Pegawai</p>
              </div>
              <div className="p-3 bg-gray-50 rounded">
                <p className="font-medium text-gray-800">Dokumen yang Dihasilkan:</p>
                <p className="mt-1">Setiap perjalanan dinas menghasilkan 2 dokumen: Nota Permintaan dan Berita Acara dengan nomor urut otomatis.</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </AdminLayout>
  );
}
