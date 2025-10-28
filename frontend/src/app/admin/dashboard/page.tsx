'use client';

import { useEffect, useState } from 'react';
import AdminLayout from '@/components/AdminLayout';
import { travelRequestAPI, employeeAPI } from '@/lib/api';
import Link from 'next/link';

interface EmployeeSPDStats {
  employee_id: number;
  employee_name: string;
  nip: string;
  position: string;
  spd_count: number;
}

export default function AdminDashboard() {
  const [stats, setStats] = useState({
    totalRequests: 0,
    totalEmployees: 0,
    totalPositions: 21,
  });
  const [employeeStats, setEmployeeStats] = useState<EmployeeSPDStats[]>([]);
  const [currentYear, setCurrentYear] = useState(new Date().getFullYear());
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadStats();
    loadEmployeeStats();
  }, [currentYear]);

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

  const loadEmployeeStats = async () => {
    setLoading(true);
    try {
      const response = await travelRequestAPI.getEmployeeStats(currentYear);
      setEmployeeStats(response.stats || []);
    } catch (error) {
      console.error('Failed to load employee stats:', error);
    } finally {
      setLoading(false);
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

  const getPodiumHeight = (rank: number) => {
    if (rank === 1) return 'h-48';
    if (rank === 2) return 'h-36';
    return 'h-28';
  };

  const getPodiumColor = (rank: number) => {
    if (rank === 1) return 'bg-gradient-to-t from-yellow-400 to-yellow-300';
    if (rank === 2) return 'bg-gradient-to-t from-gray-400 to-gray-300';
    return 'bg-gradient-to-t from-amber-600 to-amber-500';
  };

  const getMedalEmoji = (rank: number) => {
    if (rank === 1) return '🥇';
    if (rank === 2) return '🥈';
    return '🥉';
  };

  const top3 = employeeStats.slice(0, 3);
  const remaining = employeeStats.slice(3);

  // Reorder for podium display: 2nd, 1st, 3rd
  const podiumOrder = top3.length >= 2 ? [top3[1], top3[0], top3[2]].filter(Boolean) : top3;
  const podiumRanks = podiumOrder.length >= 2 ? [2, 1, 3] : [1];

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

        {/* SPD Leaderboard Section */}
        <div className="bg-white rounded-lg shadow-lg p-6 mb-8">
          <div className="flex justify-between items-center mb-6">
            <h3 className="text-xl font-bold text-gray-800">🏆 Klasemen SPD Karyawan {currentYear}</h3>
            <select
              value={currentYear}
              onChange={(e) => setCurrentYear(Number(e.target.value))}
              className="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
            >
              {Array.from({ length: 5 }, (_, i) => new Date().getFullYear() - i).map((year) => (
                <option key={year} value={year}>
                  {year}
                </option>
              ))}
            </select>
          </div>

          {loading ? (
            <div className="text-center py-12 text-gray-500">Memuat data...</div>
          ) : employeeStats.length === 0 ? (
            <div className="text-center py-12 text-gray-500">
              Belum ada data SPD untuk tahun {currentYear}
            </div>
          ) : (
            <>
              {/* Podium for Top 3 */}
              {top3.length > 0 && (
                <div className="mb-8">
                  <h4 className="text-center text-lg font-semibold text-gray-700 mb-6">
                    🎖️ Top 3 Podium
                  </h4>
                  <div className="flex items-end justify-center gap-4 mb-8">
                    {podiumOrder.map((employee, index) => {
                      const rank = podiumRanks[index];
                      return (
                        <div key={employee.employee_id} className="flex flex-col items-center">
                          {/* Employee Info */}
                          <div className="mb-2 text-center">
                            <div className="text-3xl mb-1">{getMedalEmoji(rank)}</div>
                            <div className="text-sm font-bold text-gray-800">#{rank}</div>
                            <div className="text-sm font-semibold text-gray-700 max-w-[120px] truncate">
                              {employee.employee_name}
                            </div>
                            <div className="text-xs text-gray-500">{employee.nip}</div>
                            <div className="text-lg font-bold text-blue-600 mt-1">
                              {employee.spd_count} SPD
                            </div>
                          </div>
                          {/* Podium Block */}
                          <div
                            className={`w-32 ${getPodiumHeight(rank)} ${getPodiumColor(
                              rank
                            )} rounded-t-lg flex items-center justify-center text-white font-bold text-2xl shadow-lg border-4 border-white`}
                          >
                            {rank}
                          </div>
                        </div>
                      );
                    })}
                  </div>
                </div>
              )}

              {/* Full Leaderboard Table */}
              <div className="overflow-x-auto">
                <h4 className="text-lg font-semibold text-gray-700 mb-4">📊 Semua Klasemen</h4>
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Rank
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        NIP
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Nama
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Jabatan
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        Jumlah SPD
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {employeeStats.map((employee, index) => (
                      <tr
                        key={employee.employee_id}
                        className={`${index < 3 ? 'bg-yellow-50' : 'hover:bg-gray-50'}`}
                      >
                        <td className="px-6 py-4 whitespace-nowrap">
                          <div className="flex items-center">
                            <span className="text-lg font-bold text-gray-700">#{index + 1}</span>
                            {index < 3 && (
                              <span className="ml-2 text-xl">{getMedalEmoji(index + 1)}</span>
                            )}
                          </div>
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                          {employee.nip}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                          {employee.employee_name}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                          {employee.position}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap">
                          <span className="px-3 py-1 inline-flex text-sm leading-5 font-semibold rounded-full bg-blue-100 text-blue-800">
                            {employee.spd_count} SPD
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </>
          )}
        </div>

        {/* Quick Links */}
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
      </div>
    </AdminLayout>
  );
}
