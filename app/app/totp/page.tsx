"use client";

import React, { useState } from "react";
import {
  Plus,
  Search,
  Shield,
  Clock,
  Copy,
  RefreshCw,
  QrCode,
  Settings,
  Trash2,
  Edit,
  Eye,
  EyeOff,
  X,
} from "lucide-react";

interface TOTPItem {
  id: string;
  name: string;
  issuer: string;
  secret: string;
  algorithm: string;
  digits: number;
  period: number;
  favorite: boolean;
  lastUsed?: Date;
}

export default function TOTPPage() {
  const [searchQuery, setSearchQuery] = useState("");
  const [showFavorites, setShowFavorites] = useState(false);
  const [showSecrets, setShowSecrets] = useState<Set<string>>(new Set());
  const [showAddModal, setShowAddModal] = useState(false);
  const [currentTime, setCurrentTime] = useState(Date.now());

  const mockItems: TOTPItem[] = [];

  const toggleSecretVisibility = (itemId: string) => {
    setShowSecrets((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(itemId)) {
        newSet.delete(itemId);
      } else {
        newSet.add(itemId);
      }
      return newSet;
    });
  };

  const filteredItems = mockItems.filter((item) => {
    const matchesSearch =
      item.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      item.issuer.toLowerCase().includes(searchQuery.toLowerCase());

    const matchesFavorites = !showFavorites || item.favorite;

    return matchesSearch && matchesFavorites;
  });

  const generateTOTP = (secret: string, period: number = 30) => {
    const timeSlot = Math.floor(Date.now() / 1000 / period);
    return "123456"; // Mock TOTP code
  };

  const getRemainingTime = (period: number = 30) => {
    const timeSlot = Math.floor(Date.now() / 1000 / period);
    const nextSlot = (timeSlot + 1) * period * 1000;
    return Math.floor((nextSlot - Date.now()) / 1000);
  };

  return (
    <div className="h-full flex flex-col">
      <div className="flex-shrink-0 p-6 border-b border-slate-800">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-white mb-2">
              Authentification TOTP
            </h1>
            <p className="text-slate-400">
              Gérez vos codes d'accès à usage unique
            </p>
          </div>
          <button
            onClick={() => setShowAddModal(true)}
            className="flex items-center px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
          >
            <Plus className="w-5 h-5 mr-2" />
            Ajouter TOTP
          </button>
        </div>
      </div>

      <div className="flex-1 flex overflow-hidden">
        <div className="w-64 flex-shrink-0 border-r border-slate-800 overflow-y-auto">
          <div className="p-4 space-y-6">
            <div className="bg-slate-900 rounded-lg p-4">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-sm font-semibold text-white">Filtres</h3>
                <button className="text-slate-400 hover:text-white transition-colors">
                  <RefreshCw className="w-4 h-4" />
                </button>
              </div>

              <div className="space-y-4">
                <div className="pt-4 border-t border-slate-700">
                  <button
                    onClick={() => setShowFavorites(!showFavorites)}
                    className={`
                      w-full flex items-center px-3 py-2 text-sm rounded-lg transition-colors
                      ${
                        showFavorites
                          ? "bg-yellow-600 text-white"
                          : "text-slate-300 hover:bg-slate-800 hover:text-white"
                      }
                    `}
                  >
                    <Settings className="w-4 h-4 mr-2" />
                    Favoris uniquement
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className="flex-1 flex flex-col overflow-hidden">
          <div className="flex-shrink-0 p-4 border-b border-slate-800">
            <div className="flex flex-col sm:flex-row gap-4">
              <div className="flex-1 relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-slate-400" />
                <input
                  type="text"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  placeholder="Rechercher un TOTP..."
                  className="w-full pl-10 pr-4 py-2 bg-slate-800 border border-slate-700 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
            </div>
          </div>

          <div className="flex-1 overflow-auto">
            {filteredItems.length === 0 ? (
              <div className="flex flex-col items-center justify-center h-full text-slate-400">
                <Shield className="w-16 h-16 mb-4 text-slate-600" />
                <h3 className="text-xl font-semibold mb-2">
                  Aucun TOTP configuré
                </h3>
                <p className="text-sm mb-6">
                  Commencez par ajouter votre première application TOTP
                </p>
                <button
                  onClick={() => setShowAddModal(true)}
                  className="flex items-center px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                >
                  <Plus className="w-5 h-5 mr-2" />
                  Ajouter TOTP
                </button>
              </div>
            ) : (
              <div className="p-6">
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                  {filteredItems.map((item) => {
                    const totpCode = generateTOTP(item.secret, item.period);
                    const remainingTime = getRemainingTime(item.period);
                    const showSecret = showSecrets.has(item.id);

                    return (
                      <div
                        key={item.id}
                        className="bg-slate-900 border border-slate-800 rounded-lg p-4 hover:bg-slate-800 transition-colors"
                      >
                        <div className="flex items-start justify-between mb-4">
                          <div className="flex items-center">
                            <div className="w-10 h-10 bg-blue-600 bg-opacity-20 rounded-lg flex items-center justify-center mr-3">
                              <Shield className="w-5 h-5 text-blue-500" />
                            </div>
                            <div>
                              <h3 className="text-sm font-semibold text-white">
                                {item.name}
                              </h3>
                              <p className="text-xs text-slate-400">
                                {item.issuer}
                              </p>
                            </div>
                          </div>
                          <div className="flex items-center space-x-1">
                            <button className="p-1 text-slate-400 hover:text-white transition-colors">
                              <Copy className="w-4 h-4" />
                            </button>
                            <button className="p-1 text-slate-400 hover:text-white transition-colors">
                              <Edit className="w-4 h-4" />
                            </button>
                            <button className="p-1 text-slate-400 hover:text-red-400 transition-colors">
                              <Trash2 className="w-4 h-4" />
                            </button>
                          </div>
                        </div>

                        <div className="space-y-3">
                          <div className="bg-slate-800 rounded-lg p-3">
                            <div className="flex items-center justify-between mb-2">
                              <span className="text-xs text-slate-400">
                                Code actuel
                              </span>
                              <div className="flex items-center text-xs text-slate-400">
                                <Clock className="w-3 h-3 mr-1" />
                                {remainingTime}s
                              </div>
                            </div>
                            <div className="flex items-center">
                              <span className="text-2xl font-mono font-bold text-white tracking-wider">
                                {showSecret ? totpCode : "••••••"}
                              </span>
                              <button
                                onClick={() => toggleSecretVisibility(item.id)}
                                className="ml-2 p-1 text-slate-400 hover:text-white transition-colors"
                              >
                                {showSecret ? (
                                  <EyeOff className="w-4 h-4" />
                                ) : (
                                  <Eye className="w-4 h-4" />
                                )}
                              </button>
                            </div>
                            <div className="mt-2 h-1 bg-slate-700 rounded-full overflow-hidden">
                              <div
                                className="h-full bg-blue-500 transition-all duration-1000"
                                style={{
                                  width: `${(remainingTime / item.period) * 100}%`,
                                }}
                              />
                            </div>
                          </div>

                          <div className="flex items-center justify-between text-xs text-slate-400">
                            <span>{item.digits} chiffres</span>
                            <span>{item.algorithm}</span>
                          </div>
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Modal d'ajout TOTP */}
      {showAddModal && (
        <div className="fixed inset-0 bg-black bg-opacity-30 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="bg-slate-900 border border-slate-800 rounded-lg p-6 w-full max-w-md">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-semibold text-white">Ajouter TOTP</h2>
              <button
                onClick={() => setShowAddModal(false)}
                className="text-slate-400 hover:text-white transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-4">
              <div className="text-center py-8">
                <QrCode className="w-16 h-16 mx-auto mb-4 text-slate-400" />
                <p className="text-sm text-slate-400 mb-4">
                  Scannez un code QR ou entrez manuellement la clé secrète
                </p>
              </div>

              <div className="space-y-3">
                <div>
                  <label className="block text-sm font-medium text-slate-300 mb-2">
                    Nom du compte
                  </label>
                  <input
                    type="text"
                    className="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="ex: john.doe@gmail.com"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-slate-300 mb-2">
                    Émetteur
                  </label>
                  <input
                    type="text"
                    className="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="ex: Google, GitHub..."
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-slate-300 mb-2">
                    Clé secrète
                  </label>
                  <input
                    type="text"
                    className="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Base32 secret key"
                  />
                </div>
              </div>

              <div className="pt-4 border-t border-slate-700">
                <div className="flex gap-3">
                  <button
                    onClick={() => setShowAddModal(false)}
                    className="flex-1 px-4 py-2 bg-slate-800 hover:bg-slate-700 text-white rounded-lg transition-colors"
                  >
                    Annuler
                  </button>
                  <button
                    onClick={() => {
                      console.log("Ajouter TOTP");
                      setShowAddModal(false);
                    }}
                    className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                  >
                    Ajouter
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
