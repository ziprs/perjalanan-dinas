'use client';

import { useEffect, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import { authAPI } from '@/lib/api';
import Link from 'next/link';

export default function AdminLayout({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const [username, setUsername] = useState('');

  useEffect(() => {
    if (!authAPI.isAuthenticated()) {
      router.push('/admin/login');
    } else {
      setUsername(localStorage.getItem('username') || '');
    }
  }, [router]);

  const handleLogout = () => {
    if (confirm('Apakah Anda yakin ingin logout?')) {
      authAPI.logout();
      router.push('/admin/login');
    }
  };

  const navItems = [
    { href: '/admin/dashboard', label: 'Dashboard' },
    { href: '/admin/employees', label: 'Kelola Karyawan' },
    { href: '/admin/travel-requests', label: 'Daftar Perjalanan Dinas' },
    { href: '/admin/at-cost', label: 'Klaim At-Cost' },
    { href: '/admin/monitoring-iuran', label: 'Monitoring Iuran' },
    { href: '/admin/settings', label: 'Pengaturan' },
  ];

  return (
    <div className="min-h-screen bg-gray-100">
      <nav className="bg-white shadow-lg">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16">
            <div className="flex">
              <div className="flex-shrink-0 flex items-center space-x-3">
                <img
                  src="/spd/logo-digibank.png"
                  alt="Divisi Digital Banking"
                  className="h-10 w-auto object-contain"
                />
                <h1 className="text-xl font-bold text-gray-800">Admin Panel</h1>
              </div>
              <div className="hidden sm:ml-6 sm:flex sm:space-x-8">
                {navItems.map((item) => (
                  <Link
                    key={item.href}
                    href={item.href}
                    className={`inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium ${
                      pathname === item.href
                        ? 'border-blue-500 text-gray-900'
                        : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700'
                    }`}
                  >
                    {item.label}
                  </Link>
                ))}
              </div>
            </div>
            <div className="flex items-center space-x-4">
              <span className="text-sm text-gray-600">
                Halo, <span className="font-medium">{username}</span>
              </span>
              <button
                onClick={handleLogout}
                className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 text-sm font-medium"
              >
                Logout
              </button>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        {children}
      </main>
    </div>
  );
}
