"use client";

import React, { useState } from "react";
import {
  Plus,
  Search,
  Filter,
  Key,
  CreditCard,
  User,
  FileText,
  MoreVertical,
  Star,
  Folder,
  Copy,
  Eye,
  EyeOff,
  Edit,
  Trash2,
  RefreshCw,
  X,
} from "lucide-react";

interface VaultItem {
  id: string;
  type: "login" | "card" | "identity" | "secureNote";
  name: string;
  username?: string;
  url?: string;
  favorite: boolean;
  folder?: string;
  lastModified: Date;
}

export default function VaultPage() {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedType, setSelectedType] = useState<string>("all");
  const [showFavorites, setShowFavorites] = useState(false);
  const [selectedFolder, setSelectedFolder] = useState<string>("all");
  const [showPasswords, setShowPasswords] = useState<Set<string>>(new Set());
  const [showAddItemModal, setShowAddItemModal] = useState(false);
  const [newItemType, setNewItemType] = useState<
    "login" | "card" | "identity" | "secureNote"
  >("login");

  const mockItems: VaultItem[] = [];

  const itemTypes = [
    { value: "all", label: "Tous les éléments", icon: Key },
    { value: "login", label: "Connexions", icon: Key },
    { value: "card", label: "Cartes", icon: CreditCard },
    { value: "identity", label: "Identités", icon: User },
    { value: "secureNote", label: "Notes sécurisées", icon: FileText },
  ];

  const folders = [
    { value: "all", label: "Tous les dossiers" },
    { value: "none", label: "Aucun dossier" },
  ];

  const togglePasswordVisibility = (itemId: string) => {
    setShowPasswords((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(itemId)) {
        newSet.delete(itemId);
      } else {
        newSet.add(itemId);
      }
      return newSet;
    });
  };

  const getItemIcon = (type: string) => {
    switch (type) {
      case "login":
        return Key;
      case "card":
        return CreditCard;
      case "identity":
        return User;
      case "secureNote":
        return FileText;
      default:
        return Key;
    }
  };

  const filteredItems = mockItems.filter((item) => {
    const matchesSearch =
      item.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      item.username?.toLowerCase().includes(searchQuery.toLowerCase()) ||
      item.url?.toLowerCase().includes(searchQuery.toLowerCase());

    const matchesType = selectedType === "all" || item.type === selectedType;
    const matchesFavorites = !showFavorites || item.favorite;
    const matchesFolder =
      selectedFolder === "all" ||
      (selectedFolder === "none" && !item.folder) ||
      item.folder === selectedFolder;

    return matchesSearch && matchesType && matchesFavorites && matchesFolder;
  });

  const handleAddItem = (
    type: "login" | "card" | "identity" | "secureNote",
  ) => {
    setNewItemType(type);
    setShowAddItemModal(true);
  };

  return (
    <div className="h-full flex flex-col">
      <div className="flex-shrink-0 p-6 border-b border-slate-800">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-white mb-2">Coffre</h1>
            <p className="text-slate-400">Gérez vos mots de passe et secrets</p>
          </div>
          <button
            onClick={() => setShowAddItemModal(true)}
            className="flex items-center px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
          >
            <Plus className="w-5 h-5 mr-2" />
            Ajouter un élément
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
                <div>
                  <label className="block text-xs font-medium text-slate-400 mb-2">
                    Type d'élément
                  </label>
                  <div className="space-y-1">
                    {itemTypes.map((type) => {
                      const Icon = type.icon;
                      return (
                        <button
                          key={type.value}
                          onClick={() => setSelectedType(type.value)}
                          className={`
                            w-full flex items-center px-3 py-2 text-sm rounded-lg transition-colors
                            ${
                              selectedType === type.value
                                ? "bg-blue-600 text-white"
                                : "text-slate-300 hover:bg-slate-800 hover:text-white"
                            }
                          `}
                        >
                          <Icon className="w-4 h-4 mr-2" />
                          {type.label}
                        </button>
                      );
                    })}
                  </div>
                </div>

                <div>
                  <label className="block text-xs font-medium text-slate-400 mb-2">
                    Dossiers
                  </label>
                  <div className="space-y-1">
                    {folders.map((folder) => (
                      <button
                        key={folder.value}
                        onClick={() => setSelectedFolder(folder.value)}
                        className={`
                          w-full flex items-center px-3 py-2 text-sm rounded-lg transition-colors
                          ${
                            selectedFolder === folder.value
                              ? "bg-blue-600 text-white"
                              : "text-slate-300 hover:bg-slate-800 hover:text-white"
                          }
                        `}
                      >
                        <Folder className="w-4 h-4 mr-2" />
                        {folder.label}
                      </button>
                    ))}
                  </div>
                </div>

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
                    <Star className="w-4 h-4 mr-2" />
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
                  placeholder="Rechercher dans le coffre..."
                  className="w-full pl-10 pr-4 py-2 bg-slate-800 border border-slate-700 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
              <button className="flex items-center px-4 py-2 bg-slate-800 hover:bg-slate-700 text-white rounded-lg transition-colors">
                <Filter className="w-5 h-5 mr-2" />
                Filtres
              </button>
            </div>
          </div>

          <div className="flex-1 overflow-auto">
            {filteredItems.length === 0 ? (
              <div className="flex flex-col items-center justify-center h-full text-slate-400">
                <Key className="w-16 h-16 mb-4 text-slate-600" />
                <h3 className="text-xl font-semibold mb-2">
                  Aucun élément trouvé
                </h3>
                <p className="text-sm mb-6">
                  Commencez par ajouter votre premier élément au coffre
                </p>
                <button
                  onClick={() => setShowAddItemModal(true)}
                  className="flex items-center px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                >
                  <Plus className="w-5 h-5 mr-2" />
                  Ajouter un élément
                </button>
              </div>
            ) : (
              <div className="h-full">
                <table className="w-full">
                  <thead className="bg-slate-800 sticky top-0">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                        Nom
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                        Nom d'utilisateur
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                        URL
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                        Dossier
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                        Actions
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-800">
                    {filteredItems.map((item) => {
                      const Icon = getItemIcon(item.type);
                      return (
                        <tr
                          key={item.id}
                          className="hover:bg-slate-800 transition-colors"
                        >
                          <td className="px-6 py-4 whitespace-nowrap">
                            <div className="flex items-center">
                              <Icon className="w-5 h-5 text-slate-400 mr-3" />
                              <div>
                                <div className="text-sm font-medium text-white flex items-center">
                                  {item.name}
                                  {item.favorite && (
                                    <Star className="w-4 h-4 text-yellow-500 ml-2 fill-current" />
                                  )}
                                </div>
                                <div className="text-xs text-slate-400 capitalize">
                                  {item.type}
                                </div>
                              </div>
                            </div>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            <div className="text-sm text-slate-300">
                              {item.username || "-"}
                            </div>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            <div className="text-sm text-slate-300">
                              {item.url || "-"}
                            </div>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            <div className="text-sm text-slate-300">
                              {item.folder || "-"}
                            </div>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            <div className="flex items-center space-x-2">
                              <button className="p-1 text-slate-400 hover:text-white transition-colors">
                                <Copy className="w-4 h-4" />
                              </button>
                              <button className="p-1 text-slate-400 hover:text-white transition-colors">
                                <Eye className="w-4 h-4" />
                              </button>
                              <button className="p-1 text-slate-400 hover:text-white transition-colors">
                                <Edit className="w-4 h-4" />
                              </button>
                              <button className="p-1 text-slate-400 hover:text-red-400 transition-colors">
                                <Trash2 className="w-4 h-4" />
                              </button>
                            </div>
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Modal d'ajout d'élément */}
      {showAddItemModal && (
        <div className="fixed inset-0 bg-black bg-opacity-30 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="bg-slate-900 border border-slate-800 rounded-lg p-6 w-full max-w-md">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-semibold text-white">
                Ajouter un élément
              </h2>
              <button
                onClick={() => setShowAddItemModal(false)}
                className="text-slate-400 hover:text-white transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-300 mb-3">
                  Type d'élément
                </label>
                <div className="grid grid-cols-2 gap-3">
                  <button
                    onClick={() => handleAddItem("login")}
                    className={`
                      flex flex-col items-center p-4 rounded-lg border-2 transition-colors
                      ${
                        newItemType === "login"
                          ? "border-blue-600 bg-blue-600 bg-opacity-20"
                          : "border-slate-700 hover:border-slate-600"
                      }
                    `}
                  >
                    <Key className="w-8 h-8 mb-2 text-blue-500" />
                    <span className="text-sm font-medium text-white">
                      Connexion
                    </span>
                  </button>

                  <button
                    onClick={() => handleAddItem("card")}
                    className={`
                      flex flex-col items-center p-4 rounded-lg border-2 transition-colors
                      ${
                        newItemType === "card"
                          ? "border-blue-600 bg-blue-600 bg-opacity-20"
                          : "border-slate-700 hover:border-slate-600"
                      }
                    `}
                  >
                    <CreditCard className="w-8 h-8 mb-2 text-green-500" />
                    <span className="text-sm font-medium text-white">
                      Carte
                    </span>
                  </button>

                  <button
                    onClick={() => handleAddItem("identity")}
                    className={`
                      flex flex-col items-center p-4 rounded-lg border-2 transition-colors
                      ${
                        newItemType === "identity"
                          ? "border-blue-600 bg-blue-600 bg-opacity-20"
                          : "border-slate-700 hover:border-slate-600"
                      }
                    `}
                  >
                    <User className="w-8 h-8 mb-2 text-purple-500" />
                    <span className="text-sm font-medium text-white">
                      Identité
                    </span>
                  </button>

                  <button
                    onClick={() => handleAddItem("secureNote")}
                    className={`
                      flex flex-col items-center p-4 rounded-lg border-2 transition-colors
                      ${
                        newItemType === "secureNote"
                          ? "border-blue-600 bg-blue-600 bg-opacity-20"
                          : "border-slate-700 hover:border-slate-600"
                      }
                    `}
                  >
                    <FileText className="w-8 h-8 mb-2 text-orange-500" />
                    <span className="text-sm font-medium text-white">
                      Note sécurisée
                    </span>
                  </button>
                </div>
              </div>

              <div className="pt-4 border-t border-slate-700">
                <p className="text-sm text-slate-400 mb-4">
                  Sélectionnez un type d'élément pour commencer à créer votre
                  entrée.
                </p>
                <div className="flex gap-3">
                  <button
                    onClick={() => setShowAddItemModal(false)}
                    className="flex-1 px-4 py-2 bg-slate-800 hover:bg-slate-700 text-white rounded-lg transition-colors"
                  >
                    Annuler
                  </button>
                  <button
                    onClick={() => {
                      // Ici vous pourriez rediriger vers une page de création
                      console.log(`Créer un élément de type: ${newItemType}`);
                      setShowAddItemModal(false);
                    }}
                    className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                  >
                    Continuer
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
