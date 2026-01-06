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

  const mockVaultData = {
    items: [
      {
        id: "1",
        type: "login",
        name: "Google",
        username: "user@gmail.com",
        password: "••••••••",
        url: "https://accounts.google.com",
        notes: "Compte principal",
        folder: "Personnel",
        favorite: true,
      },
      {
        id: "2",
        type: "card",
        name: "Carte Visa",
        cardNumber: "•••• •••• •••• 1234",
        cardholderName: "John Doe",
        expMonth: "12",
        expYear: "2025",
        cvv: "•••",
        brand: "Visa",
      },
      {
        id: "3",
        type: "note",
        name: "Note importante",
        notes: "Informations sensibles à conserver",
      },
    ],
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
      exportData =
        "name,username,password,url,notes\nGoogle,user@gmail.com,••••••••,https://accounts.google.com,Compte principal";
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
    if (selectedFormat === "json") {
      return JSON.stringify(mockVaultData, null, 2).substring(0, 500) + "...";
    } else if (selectedFormat === "csv") {
      return "name,username,password,url,notes\nGoogle,user@gmail.com,••••••••,https://accounts.google.com,Compte principal\n...";
    }
    return "Aperçu non disponible pour ce format";
  };

  return (
    <div className="h-full flex flex-col">
      <div className="flex-shrink-0 p-6 border-b border-slate-800">
        <div className="max-w-4xl mx-auto">
          <h1 className="text-3xl font-bold text-white mb-2">
            Exporter des données
          </h1>
          <p className="text-slate-400">
            Exportez votre coffre vers différents formats
          </p>
        </div>
      </div>

      <div className="flex-1 overflow-auto">
        <div className="max-w-4xl mx-auto p-6">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Main Content */}
            <div className="lg:col-span-2 space-y-6">
              {/* Format Selection */}
              <div className="bg-slate-900 border border-slate-800 rounded-lg p-6">
                <h2 className="text-xl font-semibold text-white mb-4">
                  Choisissez le format d'exportation
                </h2>
                <div className="space-y-3">
                  {exportFormats.map((format) => {
                    const Icon = format.icon;
                    return (
                      <button
                        key={format.id}
                        onClick={() => setSelectedFormat(format.id)}
                        className={`
                          w-full flex items-center p-4 rounded-lg border-2 transition-all
                          ${
                            selectedFormat === format.id
                              ? "border-blue-600 bg-blue-600 bg-opacity-10"
                              : "border-slate-700 hover:border-slate-600 hover:bg-slate-800"
                          }
                        `}
                      >
                        <div className="flex items-center flex-1">
                          <Icon
                            className={`w-6 h-6 mr-3 ${
                              format.encrypted
                                ? "text-green-500"
                                : "text-blue-500"
                            }`}
                          />
                          <div className="text-left flex-1">
                            <div className="flex items-center justify-between">
                              <div className="text-sm font-medium text-white">
                                {format.name}
                              </div>
                              {format.encrypted && (
                                <div className="flex items-center text-xs text-green-500">
                                  <Lock className="w-3 h-3 mr-1" />
                                  Chiffré
                                </div>
                              )}
                            </div>
                            <div className="text-xs text-slate-400 mt-1">
                              {format.description}
                            </div>
                            <div className="flex flex-wrap gap-1 mt-2">
                              {format.features.map((feature, index) => (
                                <span
                                  key={index}
                                  className="px-2 py-1 bg-slate-800 text-xs text-slate-300 rounded"
                                >
                                  {feature}
                                </span>
                              ))}
                            </div>
                          </div>
                        </div>
                      </button>
                    );
                  })}
                </div>
              </div>

              {/* Export Options */}
              <div className="bg-slate-900 border border-slate-800 rounded-lg p-6">
                <h2 className="text-xl font-semibold text-white mb-4">
                  Options d'exportation
                </h2>
                <div className="space-y-4">
                  <div className="space-y-3">
                    <label className="flex items-center space-x-3">
                      <input
                        type="checkbox"
                        checked={exportOptions.includeFolders}
                        onChange={(e) =>
                          setExportOptions({
                            ...exportOptions,
                            includeFolders: e.target.checked,
                          })
                        }
                        className="w-4 h-4 text-blue-600 bg-slate-800 border-slate-600 rounded focus:ring-blue-500"
                      />
                      <span className="text-sm text-slate-300">
                        Inclure les dossiers
                      </span>
                    </label>

                    <label className="flex items-center space-x-3">
                      <input
                        type="checkbox"
                        checked={exportOptions.includeAttachments}
                        onChange={(e) =>
                          setExportOptions({
                            ...exportOptions,
                            includeAttachments: e.target.checked,
                          })
                        }
                        className="w-4 h-4 text-blue-600 bg-slate-800 border-slate-600 rounded focus:ring-blue-500"
                      />
                      <span className="text-sm text-slate-300">
                        Inclure les pièces jointes
                      </span>
                    </label>

                    <label className="flex items-center space-x-3">
                      <input
                        type="checkbox"
                        checked={exportOptions.includeOrganizations}
                        onChange={(e) =>
                          setExportOptions({
                            ...exportOptions,
                            includeOrganizations: e.target.checked,
                          })
                        }
                        className="w-4 h-4 text-blue-600 bg-slate-800 border-slate-600 rounded focus:ring-blue-500"
                      />
                      <span className="text-sm text-slate-300">
                        Inclure les organisations
                      </span>
                    </label>

                    <label className="flex items-center space-x-3">
                      <input
                        type="checkbox"
                        checked={exportOptions.includeSends}
                        onChange={(e) =>
                          setExportOptions({
                            ...exportOptions,
                            includeSends: e.target.checked,
                          })
                        }
                        className="w-4 h-4 text-blue-600 bg-slate-800 border-slate-600 rounded focus:ring-blue-500"
                      />
                      <span className="text-sm text-slate-300">
                        Inclure les Sends
                      </span>
                    </label>
                  </div>
                </div>
              </div>

              {/* Export Actions */}
              <div className="bg-slate-900 border border-slate-800 rounded-lg p-6">
                <div className="space-y-4">
                  {!exportComplete ? (
                    <>
                      <div className="flex gap-3">
                        <button
                          onClick={() => setShowPreview(!showPreview)}
                          className="flex-1 px-4 py-2 bg-slate-800 hover:bg-slate-700 text-white rounded-lg transition-colors"
                        >
                          <Eye className="w-4 h-4 inline mr-2" />
                          Aperçu
                        </button>
                        <button
                          onClick={handleExport}
                          disabled={isExporting}
                          className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-slate-700 disabled:text-slate-500 text-white rounded-lg transition-colors"
                        >
                          {isExporting ? (
                            <div className="flex items-center justify-center">
                              <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
                              Exportation...
                            </div>
                          ) : (
                            <>
                              <Download className="w-4 h-4 inline mr-2" />
                              Exporter
                            </>
                          )}
                        </button>
                      </div>

                      {showPreview && (
                        <div className="mt-4">
                          <div className="p-4 bg-slate-800 rounded-lg">
                            <h4 className="text-sm font-medium text-white mb-2">
                              Aperçu:
                            </h4>
                            <pre className="text-xs text-slate-400 font-mono overflow-x-auto">
                              {getPreviewData()}
                            </pre>
                          </div>
                        </div>
                      )}
                    </>
                  ) : (
                    <div className="space-y-4">
                      <div className="p-4 bg-green-900 bg-opacity-20 border border-green-800 rounded-lg">
                        <div className="flex items-center">
                          <CheckCircle className="w-5 h-5 mr-3 text-green-500" />
                          <div>
                            <p className="text-sm font-medium text-white">
                              Exportation terminée!
                            </p>
                            <p className="text-xs text-slate-400">
                              Votre fichier a été téléchargé avec succès
                            </p>
                          </div>
                        </div>
                      </div>

                      <button
                        onClick={() => {
                          setExportComplete(false);
                          setShowPreview(false);
                        }}
                        className="w-full px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                      >
                        Nouvelle exportation
                      </button>
                    </div>
                  )}
                </div>
              </div>
            </div>

            {/* Sidebar */}
            <div className="space-y-6">
              {/* Export Summary */}
              <div className="bg-slate-900 border border-slate-800 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-white mb-4 flex items-center">
                  <Info className="w-5 h-5 mr-2 text-blue-500" />
                  Résumé de l'export
                </h3>
                <div className="space-y-3">
                  <div className="flex justify-between text-sm">
                    <span className="text-slate-400">Format:</span>
                    <span className="text-white">
                      {getSelectedFormat()?.name}
                    </span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-slate-400">Extension:</span>
                    <span className="text-white">
                      {getSelectedFormat()?.extension}
                    </span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-slate-400">Chiffrement:</span>
                    <span
                      className={
                        getSelectedFormat()?.encrypted
                          ? "text-green-500"
                          : "text-slate-400"
                      }
                    >
                      {getSelectedFormat()?.encrypted ? "Oui" : "Non"}
                    </span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-slate-400">Éléments:</span>
                    <span className="text-white">
                      {mockVaultData.items.length}
                    </span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-slate-400">Dossiers:</span>
                    <span className="text-white">
                      {mockVaultData.folders.length}
                    </span>
                  </div>
                </div>
              </div>

              {/* Security Warning */}
              <div className="bg-slate-900 border border-slate-800 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-white mb-4 flex items-center">
                  <AlertTriangle className="w-5 h-5 mr-2 text-yellow-500" />
                  Sécurité
                </h3>
                <div className="space-y-3 text-sm text-slate-400">
                  <p>
                    Les fichiers d'exportation contiennent des informations
                    sensibles.
                  </p>
                  <div className="space-y-2">
                    <div className="flex items-start">
                      <Lock className="w-4 h-4 mr-2 mt-0.5 text-red-500 flex-shrink-0" />
                      <span>Stockez le fichier dans un endroit sécurisé</span>
                    </div>
                    <div className="flex items-start">
                      <Lock className="w-4 h-4 mr-2 mt-0.5 text-red-500 flex-shrink-0" />
                      <span>Utilisez le format chiffré si possible</span>
                    </div>
                    <div className="flex items-start">
                      <Lock className="w-4 h-4 mr-2 mt-0.5 text-red-500 flex-shrink-0" />
                      <span>Supprimez le fichier après utilisation</span>
                    </div>
                  </div>
                </div>
              </div>

              {/* Tips */}
              <div className="bg-slate-900 border border-slate-800 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-white mb-4 flex items-center">
                  <Settings className="w-5 h-5 mr-2 text-purple-500" />
                  Conseils
                </h3>
                <div className="space-y-3 text-sm text-slate-400">
                  <p>
                    Exportez régulièrement vos données pour sauvegarder votre
                    coffre.
                  </p>
                  <p>
                    Le format JSON est recommandé pour une sauvegarde complète.
                  </p>
                  <p>
                    Le format CSV est utile pour migrer vers d'autres
                    gestionnaires.
                  </p>
                </div>
              </div>
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
