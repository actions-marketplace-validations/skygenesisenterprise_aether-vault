"use client";

import React, { useState } from "react";
import {
  Download,
  FileText,
  Shield,
  Lock,
  Unlock,
  Eye,
  EyeOff,
  CheckCircle,
  AlertTriangle,
  Info,
  Copy,
  X,
  Settings,
  Database,
  Key,
  CreditCard,
  User,
  Folder,
  Plus,
  Search,
  Filter,
  MoreVertical,
  Star,
  Edit,
  Trash2,
  RefreshCw,
} from "lucide-react";

interface ExportFormat {
  id: string;
  name: string;
  description: string;
  icon: React.ElementType;
  extension: string;
  encrypted: boolean;
  features: string[];
}

interface ExportOptions {
  format: string;
  includeAttachments: boolean;
  includeFolders: boolean;
  includeOrganizations: boolean;
  includeSends: boolean;
  password?: string;
}

interface VaultItem {
  id: string;
  type: "login" | "card" | "identity" | "secureNote";
  name: string;
  username?: string;
  url?: string;
  favorite: boolean;
  folder?: string;
  lastModified: Date;
  selected: boolean;
}

export default function ExportPage() {
  const [selectedFormat, setSelectedFormat] = useState<string>("json");
  const [exportOptions, setExportOptions] = useState<ExportOptions>({
    format: "json",
    includeAttachments: false,
    includeFolders: true,
    includeOrganizations: false,
    includeSends: false,
  });
  const [showPasswordModal, setShowPasswordModal] = useState(false);
  const [exportPassword, setExportPassword] = useState("");
  const [isExporting, setIsExporting] = useState(false);
  const [exportComplete, setExportComplete] = useState(false);
  const [showPreview, setShowPreview] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedType, setSelectedType] = useState<string>("all");
  const [showFavorites, setShowFavorites] = useState(false);
  const [selectAll, setSelectAll] = useState(false);

  const exportFormats: ExportFormat[] = [
    {
      id: "json",
      name: "Bitwarden (JSON)",
      description: "Format natif Bitwarden avec toutes les données",
      icon: Database,
      extension: ".json",
      encrypted: false,
      features: ["Mots de passe", "Notes", "Cartes", "Identités", "Dossiers"],
    },
    {
      id: "csv",
      name: "CSV",
      description: "Format universel pour les mots de passe uniquement",
      icon: FileText,
      extension: ".csv",
      encrypted: false,
      features: ["Mots de passe", "Notes"],
    },
    {
      id: "encrypted-json",
      name: "Bitwarden chiffré (JSON)",
      description: "Export chiffré avec mot de passe",
      icon: Shield,
      extension: ".json",
      encrypted: true,
      features: [
        "Toutes les données",
        "Chiffrement AES-256",
        "Protégé par mot de passe",
      ],
    },
    {
      id: "1password",
      name: "1Password (1pif)",
      description: "Format compatible 1Password",
      icon: Key,
      extension: ".1pif",
      encrypted: false,
      features: ["Mots de passe", "Notes", "Cartes", "Identités"],
    },
    {
      id: "lastpass",
      name: "LastPass (CSV)",
      description: "Format compatible LastPass",
      icon: Lock,
      extension: ".csv",
      encrypted: false,
      features: ["Mots de passe", "Notes"],
    },
  ];

  const mockVaultItems: VaultItem[] = [
    {
      id: "1",
      type: "login",
      name: "Google",
      username: "user@gmail.com",
      url: "https://accounts.google.com",
      favorite: true,
      folder: "Personnel",
      lastModified: new Date(),
      selected: true,
    },
    {
      id: "2",
      type: "card",
      name: "Carte Visa",
      favorite: false,
      folder: "Finance",
      lastModified: new Date(),
      selected: true,
    },
    {
      id: "3",
      type: "identity",
      name: "John Doe",
      favorite: false,
      folder: "Personnel",
      lastModified: new Date(),
      selected: false,
    },
    {
      id: "4",
      type: "secureNote",
      name: "Note importante",
      favorite: false,
      folder: "Travail",
      lastModified: new Date(),
      selected: true,
    },
    {
      id: "5",
      type: "login",
      name: "GitHub",
      username: "johndoe",
      url: "https://github.com",
      favorite: true,
      folder: "Développement",
      lastModified: new Date(),
      selected: true,
    },
  ];

  const [vaultItems, setVaultItems] = useState<VaultItem[]>(mockVaultItems);

  const itemTypes = [
    { value: "all", label: "Tous les éléments", icon: Key },
    { value: "login", label: "Connexions", icon: Key },
    { value: "card", label: "Cartes", icon: CreditCard },
    { value: "identity", label: "Identités", icon: User },
    { value: "secureNote", label: "Notes sécurisées", icon: FileText },
  ];

  const mockVaultData = {
    items: vaultItems.filter((item) => item.selected),
    folders: [
      { id: "1", name: "Personnel" },
      { id: "2", name: "Travail" },
    ],
    sends: [],
  };

  const handleExport = async () => {
    const format = exportFormats.find((f) => f.id === selectedFormat);

    if (format?.encrypted && !exportPassword) {
      setShowPasswordModal(true);
      return;
    }

    setIsExporting(true);

    // Simulate export process
    await new Promise((resolve) => setTimeout(resolve, 2000));

    // Generate mock export data
    let exportData = "";
    let filename = "";
    let mimeType = "";

    if (selectedFormat === "json") {
      exportData = JSON.stringify(mockVaultData, null, 2);
      filename = `aether-vault-export-${new Date().toISOString().split("T")[0]}.json`;
      mimeType = "application/json";
    } else if (selectedFormat === "csv") {
      const selectedItems = vaultItems.filter((item) => item.selected);
      exportData =
        "name,username,password,url,notes\n" +
        selectedItems
          .map(
            (item) =>
              `${item.name},${item.username || ""},••••••••,${item.url || ""},""`,
          )
          .join("\n");
      filename = `aether-vault-export-${new Date().toISOString().split("T")[0]}.csv`;
      mimeType = "text/csv";
    } else if (selectedFormat === "encrypted-json") {
      exportData = JSON.stringify(
        { encrypted: true, data: mockVaultData },
        null,
        2,
      );
      filename = `aether-vault-export-encrypted-${new Date().toISOString().split("T")[0]}.json`;
      mimeType = "application/json";
    }

    // Create download
    const blob = new Blob([exportData], { type: mimeType });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);

    setIsExporting(false);
    setExportComplete(true);
  };

  const getSelectedFormat = () => {
    return exportFormats.find((format) => format.id === selectedFormat);
  };

  const getPreviewData = () => {
    const selectedItems = vaultItems.filter((item) => item.selected);
    if (selectedFormat === "json") {
      return (
        JSON.stringify({ items: selectedItems }, null, 2).substring(0, 500) +
        "..."
      );
    } else if (selectedFormat === "csv") {
      return (
        "name,username,password,url,notes\n" +
        selectedItems
          .slice(0, 3)
          .map(
            (item) =>
              `${item.name},${item.username || ""},••••••••,${item.url || ""},""`,
          )
          .join("\n") +
        "\n..."
      );
    }
    return "Aperçu non disponible pour ce format";
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

  const toggleItemSelection = (itemId: string) => {
    setVaultItems((prev) =>
      prev.map((item) =>
        item.id === itemId ? { ...item, selected: !item.selected } : item,
      ),
    );
  };

  const toggleAllSelection = () => {
    const newSelectAll = !selectAll;
    setSelectAll(newSelectAll);
    setVaultItems((prev) =>
      prev.map((item) => ({ ...item, selected: newSelectAll })),
    );
  };

  const filteredItems = vaultItems.filter((item) => {
    const matchesSearch =
      item.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      item.username?.toLowerCase().includes(searchQuery.toLowerCase()) ||
      item.url?.toLowerCase().includes(searchQuery.toLowerCase());

    const matchesType = selectedType === "all" || item.type === selectedType;
    const matchesFavorites = !showFavorites || item.favorite;

    return matchesSearch && matchesType && matchesFavorites;
  });

  const selectedCount = vaultItems.filter((item) => item.selected).length;

  return (
    <div className="h-full flex flex-col">
      <div className="flex-shrink-0 p-6 border-b border-slate-800">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-white mb-2">Exportation</h1>
            <p className="text-slate-400">
              Exportez votre coffre vers différents formats
            </p>
          </div>
          <button
            onClick={handleExport}
            disabled={selectedCount === 0 || isExporting}
            className="flex items-center px-4 py-2 bg-slate-700 hover:bg-slate-600 disabled:bg-slate-800 disabled:text-slate-500 text-white rounded-lg transition-colors"
          >
            {isExporting ? (
              <>
                <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
                Exportation...
              </>
            ) : (
              <>
                <Download className="w-5 h-5 mr-2" />
                Exporter ({selectedCount})
              </>
            )}
          </button>
        </div>
      </div>

      <div className="flex-1 flex overflow-hidden">
        <div className="w-64 flex-shrink-0 border-r border-slate-800 overflow-y-auto">
          <div className="p-4 space-y-6">
            <div className="bg-slate-900 rounded-lg p-4">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-sm font-semibold text-white">Formats</h3>
                <button className="text-slate-400 hover:text-white transition-colors">
                  <RefreshCw className="w-4 h-4" />
                </button>
              </div>

              <div className="space-y-4">
                <div>
                  <label className="block text-xs font-medium text-slate-400 mb-2">
                    Format d'exportation
                  </label>
                  <div className="space-y-1">
                    {exportFormats.map((format) => {
                      const Icon = format.icon;
                      return (
                        <button
                          key={format.id}
                          onClick={() => setSelectedFormat(format.id)}
                          className={`
                            w-full flex items-center px-3 py-2 text-sm rounded-lg transition-colors
                            ${
                              selectedFormat === format.id
                                ? "bg-slate-700 text-white border border-slate-600"
                                : "text-slate-300 hover:bg-slate-800 hover:text-white"
                            }
                          `}
                        >
                          <Icon
                            className={`w-4 h-4 mr-2 ${selectedFormat === format.id ? "text-white" : format.encrypted ? "text-green-500" : "text-blue-500"}`}
                          />
                          <div className="flex-1 text-left">
                            <div>{format.name}</div>
                            {format.encrypted && (
                              <div className="text-xs opacity-75">Chiffré</div>
                            )}
                          </div>
                        </button>
                      );
                    })}
                  </div>
                </div>

                <div className="pt-4 border-t border-slate-700">
                  <label className="block text-xs font-medium text-slate-400 mb-2">
                    Type d'éléments
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

            {/* Export Options */}
            <div className="bg-slate-900 rounded-lg p-4">
              <h3 className="text-sm font-semibold text-white mb-4">Options</h3>
              <div className="space-y-3">
                <label className="flex items-center space-x-2 text-sm text-slate-300">
                  <input
                    type="checkbox"
                    checked={exportOptions.includeFolders}
                    onChange={(e) =>
                      setExportOptions({
                        ...exportOptions,
                        includeFolders: e.target.checked,
                      })
                    }
                    className="w-3 h-3 text-blue-600 bg-slate-800 border-slate-600 rounded"
                  />
                  Inclure les dossiers
                </label>

                <label className="flex items-center space-x-2 text-sm text-slate-300">
                  <input
                    type="checkbox"
                    checked={exportOptions.includeAttachments}
                    onChange={(e) =>
                      setExportOptions({
                        ...exportOptions,
                        includeAttachments: e.target.checked,
                      })
                    }
                    className="w-3 h-3 text-blue-600 bg-slate-800 border-slate-600 rounded"
                  />
                  Inclure les pièces jointes
                </label>

                <label className="flex items-center space-x-2 text-sm text-slate-300">
                  <input
                    type="checkbox"
                    checked={exportOptions.includeOrganizations}
                    onChange={(e) =>
                      setExportOptions({
                        ...exportOptions,
                        includeOrganizations: e.target.checked,
                      })
                    }
                    className="w-3 h-3 text-blue-600 bg-slate-800 border-slate-600 rounded"
                  />
                  Inclure les organisations
                </label>

                <label className="flex items-center space-x-2 text-sm text-slate-300">
                  <input
                    type="checkbox"
                    checked={exportOptions.includeSends}
                    onChange={(e) =>
                      setExportOptions({
                        ...exportOptions,
                        includeSends: e.target.checked,
                      })
                    }
                    className="w-3 h-3 text-blue-600 bg-slate-800 border-slate-600 rounded"
                  />
                  Inclure les Sends
                </label>
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
                  placeholder="Rechercher dans les éléments à exporter..."
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
            <div className="p-6 space-y-6">
              {/* Export Summary */}
              <div className="bg-slate-900 border border-slate-800 rounded-xl p-6">
                <h2 className="text-lg font-semibold text-white mb-4">
                  Résumé de l'exportation
                </h2>
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                  <div className="text-center p-3 bg-slate-800 rounded-lg">
                    <div className="text-xl font-bold text-blue-500">
                      {selectedCount}
                    </div>
                    <div className="text-xs text-slate-400">Sélectionnés</div>
                  </div>
                  <div className="text-center p-3 bg-slate-800 rounded-lg">
                    <div className="text-xl font-bold text-white">
                      {vaultItems.length}
                    </div>
                    <div className="text-xs text-slate-400">Total</div>
                  </div>
                  <div className="text-center p-3 bg-slate-800 rounded-lg">
                    <div className="text-xl font-bold text-green-500">
                      {getSelectedFormat()?.name.split(" ")[0]}
                    </div>
                    <div className="text-xs text-slate-400">Format</div>
                  </div>
                  <div className="text-center p-3 bg-slate-800 rounded-lg">
                    <div className="text-xl font-bold text-purple-500">
                      {getSelectedFormat()?.encrypted ? "Oui" : "Non"}
                    </div>
                    <div className="text-xs text-slate-400">Chiffré</div>
                  </div>
                </div>
              </div>

              {/* Items Table */}
              <div className="bg-slate-900 border border-slate-800 rounded-xl">
                <table className="w-full">
                  <thead className="bg-slate-800 sticky top-0">
                    <tr>
                      <th className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                        <input
                          type="checkbox"
                          checked={selectAll}
                          onChange={toggleAllSelection}
                          className="w-3 h-3 text-blue-600 bg-slate-800 border-slate-600 rounded"
                        />
                      </th>
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
                          className={`hover:bg-slate-800 transition-colors ${
                            item.selected ? "bg-slate-700 bg-opacity-30" : ""
                          }`}
                        >
                          <td className="px-6 py-4 whitespace-nowrap">
                            <input
                              type="checkbox"
                              checked={item.selected}
                              onChange={() => toggleItemSelection(item.id)}
                              className="w-3 h-3 text-blue-600 bg-slate-800 border-slate-600 rounded"
                            />
                          </td>
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
                                  {item.type === "login" && "Connexion"}
                                  {item.type === "card" && "Carte"}
                                  {item.type === "identity" && "Identité"}
                                  {item.type === "secureNote" &&
                                    "Note sécurisée"}
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
                            </div>
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </table>
              </div>

              {/* Preview */}
              {showPreview && selectedCount > 0 && (
                <div className="bg-slate-900 border border-slate-800 rounded-xl p-6">
                  <h3 className="text-lg font-semibold text-white mb-4">
                    Aperçu de l'exportation
                  </h3>
                  <div className="p-4 bg-slate-800 rounded-lg">
                    <pre className="text-xs text-slate-400 font-mono overflow-x-auto">
                      {getPreviewData()}
                    </pre>
                  </div>
                </div>
              )}

              {/* Empty State */}
              {filteredItems.length === 0 && (
                <div className="flex flex-col items-center justify-center text-slate-400">
                  <Database className="w-16 h-16 mb-4 text-slate-600" />
                  <h3 className="text-xl font-semibold mb-2">
                    Aucun élément à exporter
                  </h3>
                  <p className="text-sm mb-6">
                    Aucun élément ne correspond à vos critères de recherche
                  </p>
                  <button
                    onClick={() => {
                      setSearchQuery("");
                      setSelectedType("all");
                      setShowFavorites(false);
                    }}
                    className="flex items-center px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                  >
                    <RefreshCw className="w-5 h-5 mr-2" />
                    Réinitialiser les filtres
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Password Modal */}
      {showPasswordModal && (
        <div className="fixed inset-0 bg-black bg-opacity-30 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="bg-slate-900 border border-slate-800 rounded-lg p-6 w-full max-w-md">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-semibold text-white">
                Mot de passe d'exportation
              </h2>
              <button
                onClick={() => setShowPasswordModal(false)}
                className="text-slate-400 hover:text-white transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-slate-300 mb-2">
                  Entrez un mot de passe pour chiffrer l'export
                </label>
                <input
                  type="password"
                  value={exportPassword}
                  onChange={(e) => setExportPassword(e.target.value)}
                  placeholder="Mot de passe fort"
                  className="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>

              <div className="pt-4 border-t border-slate-700">
                <div className="flex gap-3">
                  <button
                    onClick={() => setShowPasswordModal(false)}
                    className="flex-1 px-4 py-2 bg-slate-800 hover:bg-slate-700 text-white rounded-lg transition-colors"
                  >
                    Annuler
                  </button>
                  <button
                    onClick={() => {
                      setShowPasswordModal(false);
                      setExportOptions({
                        ...exportOptions,
                        password: exportPassword,
                      });
                      handleExport();
                    }}
                    disabled={!exportPassword}
                    className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-slate-700 disabled:text-slate-500 text-white rounded-lg transition-colors"
                  >
                    Exporter
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
