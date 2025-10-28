'use client';

import { useState } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import SPDForm from '@/components/SPDForm';
import AtCostForm from '@/components/AtCostForm';

export default function Home() {
  const [activeTab, setActiveTab] = useState<'spd' | 'atcost'>('spd');

  return (
    <main className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 py-12 px-4">
      <div className="max-w-4xl mx-auto">
        <div className="bg-white rounded-lg shadow-xl p-8">
          {/* Logo Section */}
          <div className="flex justify-center mb-6">
            <Image
              src="/spd/logo-digibank.png"
              alt="Divisi Digital Banking"
              width={200}
              height={96}
              className="h-24 w-auto object-contain"
              priority
              unoptimized
            />
          </div>

          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-gray-800 mb-2">
              SISTEM PERJALANAN DINAS
            </h1>
            <p className="text-gray-600">
              Form Pengajuan SPD dan Klaim At-Cost
            </p>
          </div>

          {/* Tab Navigation */}
          <div className="flex border-b border-gray-200 mb-6">
            <button
              onClick={() => setActiveTab('spd')}
              className={`flex-1 py-3 px-4 text-center font-medium transition-colors ${
                activeTab === 'spd'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-800'
              }`}
            >
              Buat Nota Perjalanan Dinas
            </button>
            <button
              onClick={() => setActiveTab('atcost')}
              className={`flex-1 py-3 px-4 text-center font-medium transition-colors ${
                activeTab === 'atcost'
                  ? 'border-b-2 border-blue-600 text-blue-600'
                  : 'text-gray-600 hover:text-gray-800'
              }`}
            >
              Buat Klaim At-Cost
            </button>
          </div>

          {/* Form Content */}
          {activeTab === 'spd' ? <SPDForm /> : <AtCostForm />}

          {/* Admin Login Link */}
          <div className="mt-6 text-center">
            <Link
              href="/admin/login"
              className="inline-block px-6 py-3 bg-gray-600 text-white rounded-lg hover:bg-gray-700 font-medium transition-colors"
            >
              Admin Login
            </Link>
          </div>
        </div>

        <div className="mt-6 text-center text-gray-600 text-sm">
          <p>UncleJSS</p>
        </div>
      </div>
    </main>
  );
}
