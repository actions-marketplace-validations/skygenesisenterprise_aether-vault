"use client";

import React, { useState } from "react";
import {
  Plus,
  Search,
  Link,
  FileText,
  Clock,
  Eye,
  EyeOff,
  Copy,
  Trash2,
  Edit,
  Shield,
  Calendar,
  Users,
  X,
  Upload,
  Lock,
  Unlock,
  CheckCircle,
  AlertCircle,
} from "lucide-react";

interface SendItem {
  id: string;
  type: "text" | "file";
  name: string;
  content?: string;
  fileName?: string;
  fileSize?: number;
  password?: string;
  accessCount: number;
  maxAccessCount: number;
  createdAt: Date;
  expiresAt: Date;
  hasPassword: boolean;
  disabled: boolean;
}

export default function SendsPage() {
  const [searchQuery, setSearchQuery] = useState("");
  const [showDisabled, setShowDisabled] = useState(false);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [newSendType, setNewSendType] = useState<"text" | "file">("text");
  const [showPasswords, setShowPasswords] = useState<Set<string>>(new Set());

  const mockSends: SendItem[] = [];

  const togglePasswordVisibility = (sendId: string) => {
    setShowPasswords((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(sendId)) {
        newSet.delete(sendId);
      } else {
        newSet.add(sendId);
      }
      return newSet;
    });
  };

  const filteredSends = mockSends.filter((send) => {
    const matchesSearch =
      send.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      send.fileName?.toLowerCase().includes(searchQuery.toLowerCase());

    const matchesDisabled = !showDisabled || !send.disabled;

    return matchesSearch && matchesDisabled;
  });

  const getStatusColor = (send: SendItem) => {
    if (send.disabled) return "text-slate-500";
    if (send.accessCount >= send.maxAccessCount) return "text-red-500";
    if (new Date() > send.expiresAt) return "text-orange-500";
    return "text-green-500";
  };

  const getStatusIcon = (send: SendItem) => {
    if (send.disabled) return <Lock className="w-4 h-4" />;
    if (send.accessCount >= send.maxAccessCount)
      return <AlertCircle className="w-4 h-4" />;
    if (new Date() > send.expiresAt) return <Clock className="w-4 h-4" />;
    return <CheckCircle className="w-4 h-4" />;
  };

  const getStatusText = (send: SendItem) => {
    if (send.disabled) return "Désactivé";
    if (send.accessCount >= send.maxAccessCount) return "Limite atteinte";
    if (new Date() > send.expiresAt) return "Expiré";
    return "Actif";
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const sizes = ["Bytes", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
  };

  const getRemainingTime = (expiresAt: Date) => {
    const now = new Date();
    const diff = expiresAt.getTime() - now.getTime();

    if (diff <= 0) return "Expiré";

    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));

    if (days > 0) return `${days}j ${hours}h`;
    if (hours > 0) return `${hours}h`;
    return "Moins d'1h";
  };

  return (
    <div className="h-full flex flex-col">
      <div className="flex-shrink-0 p-6 border-b border-slate-800">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-white mb-2">Sends</h1>
            <p className="text-slate-400">
              Partagez des informations sensibles de manière sécurisée
            </p>
          </div>
          <button
            onClick={() => setShowCreateModal(true)}
            className="flex items-center px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
          >
            <Plus className="w-5 h-5 mr-2" />
            Nouveau Send
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
                  <Clock className="w-4 h-4" />
                </button>
              </div>

              <div className="space-y-4">
                <div>
                  <label className="block text-xs font-medium text-slate-400 mb-2">
                    Type
                  </label>
                  <div className="space-y-1">
                    <button className="w-full flex items-center px-3 py-2 text-sm rounded-lg bg-blue-600 text-white">
                      <FileText className="w-4 h-4 mr-2" />
                      Tous les types
                    </button>
                    <button className="w-full flex items-center px-3 py-2 text-sm rounded-lg text-slate-300 hover:bg-slate-800 hover:text-white transition-colors">
                      <FileText className="w-4 h-4 mr-2" />
                      Texte
                    </button>
                    <button className="w-full flex items-center px-3 py-2 text-sm rounded-lg text-slate-300 hover:bg-slate-800 hover:text-white transition-colors">
                      <Upload className="w-4 h-4 mr-2" />
                      Fichier
                    </button>
                  </div>
                </div>

                <div className="pt-4 border-t border-slate-700">
                  <button
                    onClick={() => setShowDisabled(!showDisabled)}
                    className={`
                      w-full flex items-center px-3 py-2 text-sm rounded-lg transition-colors
                      ${
                        showDisabled
                          ? "bg-slate-700 text-white"
                          : "text-slate-300 hover:bg-slate-800 hover:text-white"
                      }
                    `}
                  >
                    <Lock className="w-4 h-4 mr-2" />
                    Afficher désactivés
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
                  placeholder="Rechercher un Send..."
                  className="w-full pl-10 pr-4 py-2 bg-slate-800 border border-slate-700 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
            </div>
          </div>

          <div className="flex-1 overflow-auto">
            {filteredSends.length === 0 ? (
              <div className="flex flex-col items-center justify-center h-full text-slate-400">
                <Shield className="w-16 h-16 mb-4 text-slate-600" />
                <h3 className="text-xl font-semibold mb-2">Aucun Send créé</h3>
                <p className="text-sm mb-6">
                  Commencez par créer votre premier Send pour partager des
                  informations
                </p>
                <button
                  onClick={() => setShowCreateModal(true)}
                  className="flex items-center px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                >
                  <Plus className="w-5 h-5 mr-2" />
                  Nouveau Send
                </button>
              </div>
            ) : (
              <div className="p-6">
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                  {filteredSends.map((send) => (
                    <div
                      key={send.id}
                      className="bg-slate-900 border border-slate-800 rounded-lg p-4 hover:bg-slate-800 transition-colors"
                    >
                      <div className="flex items-start justify-between mb-4">
                        <div className="flex items-center">
                          <div className="w-10 h-10 bg-blue-600 bg-opacity-20 rounded-lg flex items-center justify-center mr-3">
                            {send.type === "file" ? (
                              <Upload className="w-5 h-5 text-blue-500" />
                            ) : (
                              <FileText className="w-5 h-5 text-blue-500" />
                            )}
                          </div>
                          <div>
                            <h3 className="text-sm font-semibold text-white">
                              {send.name}
                            </h3>
                            <div
                              className={`flex items-center text-xs ${getStatusColor(send)}`}
                            >
                              {getStatusIcon(send)}
                              <span className="ml-1">
                                {getStatusText(send)}
                              </span>
                            </div>
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
                        <div className="flex items-center justify-between text-xs text-slate-400">
                          <div className="flex items-center">
                            <Users className="w-3 h-3 mr-1" />
                            <span>
                              {send.accessCount}/{send.maxAccessCount}
                            </span>
                          </div>
                          <div className="flex items-center">
                            <Clock className="w-3 h-3 mr-1" />
                            <span>{getRemainingTime(send.expiresAt)}</span>
                          </div>
                        </div>

                        {send.type === "file" && send.fileSize && (
                          <div className="text-xs text-slate-400">
                            Taille: {formatFileSize(send.fileSize)}
                          </div>
                        )}

                        {send.hasPassword && (
                          <div className="flex items-center text-xs text-slate-400">
                            <Lock className="w-3 h-3 mr-1" />
                            <span>Protégé par mot de passe</span>
                          </div>
                        )}

                        <div className="pt-3 border-t border-slate-700">
                          <button className="w-full flex items-center justify-center px-3 py-2 bg-slate-800 hover:bg-slate-700 text-white text-sm rounded-lg transition-colors">
                            <Link className="w-4 h-4 mr-2" />
                            Copier le lien
                          </button>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Modal de création de Send */}
      {showCreateModal && (
        <div className="fixed inset-0 bg-black bg-opacity-30 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="bg-slate-900 border border-slate-800 rounded-lg p-6 w-full max-w-2xl">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-semibold text-white">Nouveau Send</h2>
              <button
                onClick={() => setShowCreateModal(false)}
                className="text-slate-400 hover:text-white transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-6">
              <div>
                <label className="block text-sm font-medium text-slate-300 mb-3">
                  Type de Send
                </label>
                <div className="grid grid-cols-2 gap-4">
                  <button
                    onClick={() => setNewSendType("text")}
                    className={`
                      flex flex-col items-center p-4 rounded-lg border-2 transition-colors
                      ${
                        newSendType === "text"
                          ? "border-blue-600 bg-blue-600 bg-opacity-20"
                          : "border-slate-700 hover:border-slate-600"
                      }
                    `}
                  >
                    <FileText className="w-8 h-8 mb-2 text-blue-500" />
                    <span className="text-sm font-medium text-white">
                      Texte
                    </span>
                  </button>

                  <button
                    onClick={() => setNewSendType("file")}
                    className={`
                      flex flex-col items-center p-4 rounded-lg border-2 transition-colors
                      ${
                        newSendType === "file"
                          ? "border-blue-600 bg-blue-600 bg-opacity-20"
                          : "border-slate-700 hover:border-slate-600"
                      }
                    `}
                  >
                    <Upload className="w-8 h-8 mb-2 text-green-500" />
                    <span className="text-sm font-medium text-white">
                      Fichier
                    </span>
                  </button>
                </div>
              </div>

              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-slate-300 mb-2">
                    Nom
                  </label>
                  <input
                    type="text"
                    className="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Nom du Send"
                  />
                </div>

                {newSendType === "text" ? (
                  <div>
                    <label className="block text-sm font-medium text-slate-300 mb-2">
                      Texte à partager
                    </label>
                    <textarea
                      className="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent h-32 resize-none"
                      placeholder="Entrez le texte à partager..."
                    />
                  </div>
                ) : (
                  <div>
                    <label className="block text-sm font-medium text-slate-300 mb-2">
                      Fichier à partager
                    </label>
                    <div className="border-2 border-dashed border-slate-700 rounded-lg p-8 text-center hover:border-slate-600 transition-colors cursor-pointer">
                      <Upload className="w-12 h-12 mx-auto mb-3 text-slate-400" />
                      <p className="text-sm text-slate-400">
                        Cliquez pour sélectionner un fichier ou glissez-déposez
                      </p>
                      <p className="text-xs text-slate-500 mt-2">
                        Taille maximale: 500 MB
                      </p>
                    </div>
                  </div>
                )}

                <div>
                  <label className="block text-sm font-medium text-slate-300 mb-2">
                    Mot de passe (optionnel)
                  </label>
                  <input
                    type="password"
                    className="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    placeholder="Laisser vide pour un accès sans mot de passe"
                  />
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-slate-300 mb-2">
                      Date d'expiration
                    </label>
                    <select className="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent">
                      <option value="1">1 jour</option>
                      <option value="7">7 jours</option>
                      <option value="30">30 jours</option>
                      <option value="custom">Personnalisé</option>
                    </select>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-slate-300 mb-2">
                      Nombre d'accès maximal
                    </label>
                    <select className="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent">
                      <option value="1">1 accès</option>
                      <option value="3">3 accès</option>
                      <option value="5">5 accès</option>
                      <option value="-1">Illimité</option>
                    </select>
                  </div>
                </div>
              </div>

              <div className="pt-4 border-t border-slate-700">
                <div className="flex gap-3">
                  <button
                    onClick={() => setShowCreateModal(false)}
                    className="flex-1 px-4 py-2 bg-slate-800 hover:bg-slate-700 text-white rounded-lg transition-colors"
                  >
                    Annuler
                  </button>
                  <button
                    onClick={() => {
                      console.log("Créer Send");
                      setShowCreateModal(false);
                    }}
                    className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                  >
                    Créer le Send
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
