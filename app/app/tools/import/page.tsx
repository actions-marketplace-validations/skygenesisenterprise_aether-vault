"use client";

import React, { useState } from "react";
import {
  Upload,
  FileText,
  AlertTriangle,
  CheckCircle,
  Info,
  Download,
  Folder,
  Lock,
  Unlock,
  X,
  ArrowRight,
  Database,
  Shield,
  Plus,
  Search,
  Filter,
  Key,
  CreditCard,
  User,
  MoreVertical,
  Star,
  Copy,
  Eye,
  EyeOff,
  Edit,
  Trash2,
  RefreshCw,
  Clock,
} from "lucide-react";

interface ImportSource {
  id: string;
  name: string;
  description: string;
  icon: React.ElementType;
  formats: string[];
  color: string;
}

interface ImportStep {
  id: number;
  title: string;
  description: string;
  completed: boolean;
}

interface ImportItem {
  id: string;
  name: string;
  type: string;
  source: string;
  status: "pending" | "success" | "failed" | "skipped";
  timestamp: Date;
}

export default function ImportPage() {
  const [selectedSource, setSelectedSource] = useState<string>("");
  const [importFile, setImportFile] = useState<File | null>(null);
  const [importStep, setImportStep] = useState(1);
  const [isImporting, setIsImporting] = useState(false);
  const [importComplete, setImportComplete] = useState(false);
  const [importResults, setImportResults] = useState({
    successful: 0,
    failed: 0,
    skipped: 0,
  });
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedType, setSelectedType] = useState<string>("all");
  const [showFavorites, setShowFavorites] = useState(false);
  const [importItems, setImportItems] = useState<ImportItem[]>([]);

  const importSources: ImportSource[] = [
    {
      id: "lastpass",
      name: "LastPass",
      description: "Importez depuis LastPass (CSV ou JSON)",
      icon: Lock,
      formats: [".csv", ".json"],
      color: "text-red-500",
    },
    {
      id: "1password",
      name: "1Password",
      description: "Importez depuis 1Password (1pif, CSV)",
      icon: Shield,
      formats: [".1pif", ".csv"],
      color: "text-blue-500",
    },
    {
      id: "chrome",
      name: "Google Chrome",
      description: "Importez depuis Chrome (CSV)",
      icon: Database,
      formats: [".csv"],
      color: "text-green-500",
    },
    {
      id: "firefox",
      name: "Mozilla Firefox",
      description: "Importez depuis Firefox (CSV)",
      icon: Database,
      formats: [".csv"],
      color: "text-orange-500",
    },
    {
      id: "dashlane",
      name: "Dashlane",
      description: "Importez depuis Dashlane (CSV ou JSON)",
      icon: Lock,
      formats: [".csv", ".json"],
      color: "text-purple-500",
    },
    {
      id: "keepass",
      name: "KeePass",
      description: "Importez depuis KeePass (XML)",
      icon: Shield,
      formats: [".xml"],
      color: "text-teal-500",
    },
    {
      id: "bitwarden",
      name: "Bitwarden",
      description: "Importez depuis un autre coffre Bitwarden (JSON)",
      icon: Shield,
      formats: [".json"],
      color: "text-blue-600",
    },
    {
      id: "csv",
      name: "Fichier CSV",
      description: "Importez depuis un fichier CSV générique",
      icon: FileText,
      formats: [".csv"],
      color: "text-slate-500",
    },
  ];

  const itemTypes = [
    { value: "all", label: "Tous les éléments", icon: Key },
    { value: "login", label: "Connexions", icon: Key },
    { value: "card", label: "Cartes", icon: CreditCard },
    { value: "identity", label: "Identités", icon: User },
    { value: "secureNote", label: "Notes sécurisées", icon: FileText },
  ];

  const importSteps: ImportStep[] = [
    {
      id: 1,
      title: "Choisir la source",
      description: "Sélectionnez d'où vous importez vos données",
      completed: selectedSource !== "",
    },
    {
      id: 2,
      title: "Télécharger le fichier",
      description: "Uploadez votre fichier d'exportation",
      completed: importFile !== null,
    },
    {
      id: 3,
      title: "Vérifier et importer",
      description: "Confirmez les détails et lancez l'importation",
      completed: importComplete,
    },
  ];

  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setImportFile(file);
      setImportStep(3);
      // Generate mock import items
      generateMockItems();
    }
  };

  const generateMockItems = () => {
    const mockItems: ImportItem[] = [];
    const types = ["login", "card", "identity", "secureNote"];
    const statuses: ("pending" | "success" | "failed" | "skipped")[] = [
      "success",
      "failed",
      "skipped",
    ];

    for (let i = 0; i < 20; i++) {
      mockItems.push({
        id: `item-${i}`,
        name: `Élément d'importation ${i + 1}`,
        type: types[Math.floor(Math.random() * types.length)],
        source: selectedSource,
        status:
          i < 15
            ? "success"
            : statuses[Math.floor(Math.random() * statuses.length)],
        timestamp: new Date(),
      });
    }
    setImportItems(mockItems);
  };

  const handleImport = async () => {
    if (!importFile || !selectedSource) return;

    setIsImporting(true);

    // Simulate import process
    await new Promise((resolve) => setTimeout(resolve, 2000));

    // Mock results
    const results = {
      successful: Math.floor(Math.random() * 50) + 10,
      failed: Math.floor(Math.random() * 5),
      skipped: Math.floor(Math.random() * 3),
    };
    setImportResults(results);

    // Update import items status
    setImportItems((prev) =>
      prev.map((item) => ({
        ...item,
        status:
          Math.random() > 0.2
            ? "success"
            : Math.random() > 0.5
              ? "failed"
              : "skipped",
      })),
    );

    setIsImporting(false);
    setImportComplete(true);
  };

  const resetImport = () => {
    setSelectedSource("");
    setImportFile(null);
    setImportStep(1);
    setIsImporting(false);
    setImportComplete(false);
    setImportResults({ successful: 0, failed: 0, skipped: 0 });
    setImportItems([]);
  };

  const getSelectedSource = () => {
    return importSources.find((source) => source.id === selectedSource);
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

  const getStatusIcon = (status: string) => {
    switch (status) {
      case "success":
        return CheckCircle;
      case "failed":
        return AlertTriangle;
      case "skipped":
        return AlertTriangle;
      default:
        return Clock;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "success":
        return "text-green-500";
      case "failed":
        return "text-red-500";
      case "skipped":
        return "text-yellow-500";
      default:
        return "text-slate-400";
    }
  };

  const filteredItems = importItems.filter((item) => {
    const matchesSearch =
      item.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      item.type.toLowerCase().includes(searchQuery.toLowerCase()) ||
      item.source.toLowerCase().includes(searchQuery.toLowerCase());

    const matchesType = selectedType === "all" || item.type === selectedType;
    const matchesFavorites = !showFavorites || item.status === "success";

    return matchesSearch && matchesType && matchesFavorites;
  });

  return (
    <div className="h-full flex flex-col">
      <div className="flex-shrink-0 p-6 border-b border-slate-800">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-white mb-2">Importation</h1>
            <p className="text-slate-400">
              Importez vos mots de passe et données depuis d'autres
              gestionnaires
            </p>
          </div>
          <button
            onClick={resetImport}
            className="flex items-center px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
          >
            <Plus className="w-5 h-5 mr-2" />
            Nouvelle importation
          </button>
        </div>
      </div>

      <div className="flex-1 flex overflow-hidden">
        <div className="w-64 flex-shrink-0 border-r border-slate-800 overflow-y-auto">
          <div className="p-4 space-y-6">
            <div className="bg-slate-900 rounded-lg p-4">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-sm font-semibold text-white">Sources</h3>
                <button className="text-slate-400 hover:text-white transition-colors">
                  <RefreshCw className="w-4 h-4" />
                </button>
              </div>

              <div className="space-y-4">
                <div>
                  <label className="block text-xs font-medium text-slate-400 mb-2">
                    Source d'importation
                  </label>
                  <div className="space-y-1">
                    {importSources.map((source) => {
                      const Icon = source.icon;
                      return (
                        <button
                          key={source.id}
                          onClick={() => {
                            setSelectedSource(source.id);
                            setImportStep(2);
                          }}
                          className={`
                            w-full flex items-center px-3 py-2 text-sm rounded-lg transition-colors
                            ${
                              selectedSource === source.id
                                ? "bg-blue-600 text-white"
                                : "text-slate-300 hover:bg-slate-800 hover:text-white"
                            }
                          `}
                        >
                          <Icon
                            className={`w-4 h-4 mr-2 ${selectedSource === source.id ? "text-white" : source.color}`}
                          />
                          {source.name}
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
                    Réussis uniquement
                  </button>
                </div>
              </div>
            </div>

            {/* Import Progress */}
            {selectedSource && (
              <div className="bg-slate-900 rounded-lg p-4">
                <h3 className="text-sm font-semibold text-white mb-4">
                  Progression
                </h3>
                <div className="space-y-3">
                  {importSteps.map((step) => (
                    <div key={step.id} className="flex items-center">
                      <div
                        className={`
                          w-6 h-6 rounded-full flex items-center justify-center text-xs font-medium mr-3
                          ${
                            step.completed
                              ? "bg-green-600 text-white"
                              : importStep >= step.id
                                ? "bg-slate-700 text-white"
                                : "bg-slate-800 text-slate-400"
                          }
                        `}
                      >
                        {step.completed ? (
                          <CheckCircle className="w-3 h-3" />
                        ) : (
                          step.id
                        )}
                      </div>
                      <div className="flex-1">
                        <div
                          className={`text-xs font-medium ${step.completed ? "text-white" : "text-slate-400"}`}
                        >
                          {step.title}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
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
                  placeholder="Rechercher dans les éléments importés..."
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
              {/* Import Steps Content */}
              {!importComplete && (
                <div className="bg-slate-900 border border-slate-800 rounded-xl p-6">
                  {/* Step 1: Select Source */}
                  {importStep === 1 && (
                    <div>
                      <h2 className="text-xl font-semibold text-white mb-4">
                        Choisissez la source d'importation
                      </h2>
                      <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                        {importSources.map((source) => {
                          const Icon = source.icon;
                          return (
                            <button
                              key={source.id}
                              onClick={() => {
                                setSelectedSource(source.id);
                                setImportStep(2);
                              }}
                              className={`
                                flex items-center p-4 rounded-lg border-2 transition-all
                                ${
                                  selectedSource === source.id
                                    ? "border-blue-600 bg-blue-600 bg-opacity-10"
                                    : "border-slate-700 hover:border-slate-600 hover:bg-slate-800"
                                }
                              `}
                            >
                              <Icon
                                className={`w-6 h-6 mr-3 ${source.color}`}
                              />
                              <div className="text-left">
                                <div className="text-sm font-medium text-white">
                                  {source.name}
                                </div>
                                <div className="text-xs text-slate-400">
                                  {source.formats.join(", ")}
                                </div>
                              </div>
                            </button>
                          );
                        })}
                      </div>
                    </div>
                  )}

                  {/* Step 2: Upload File */}
                  {importStep === 2 && selectedSource && (
                    <div>
                      <h2 className="text-xl font-semibold text-white mb-4">
                        Téléchargez votre fichier
                      </h2>

                      <div className="mb-4">
                        <div className="flex items-center p-3 bg-slate-800 rounded-lg">
                          {(() => {
                            const Icon = getSelectedSource()?.icon || FileText;
                            const color =
                              getSelectedSource()?.color || "text-slate-500";
                            return (
                              <>
                                <Icon className={`w-5 h-5 mr-3 ${color}`} />
                                <div>
                                  <div className="text-sm font-medium text-white">
                                    {getSelectedSource()?.name}
                                  </div>
                                  <div className="text-xs text-slate-400">
                                    Formats supportés:{" "}
                                    {getSelectedSource()?.formats.join(", ")}
                                  </div>
                                </div>
                              </>
                            );
                          })()}
                        </div>
                      </div>

                      <div className="border-2 border-dashed border-slate-700 rounded-lg p-8 text-center hover:border-slate-600 transition-colors">
                        <Upload className="w-12 h-12 mx-auto mb-4 text-slate-400" />
                        <label className="cursor-pointer">
                          <span className="text-sm text-slate-300 mb-2 block">
                            Cliquez pour sélectionner un fichier ou
                            glissez-déposez
                          </span>
                          <span className="text-xs text-slate-500">
                            Taille maximale: 50 MB
                          </span>
                          <input
                            type="file"
                            accept={getSelectedSource()?.formats.join(",")}
                            onChange={handleFileUpload}
                            className="hidden"
                          />
                        </label>
                      </div>

                      {importFile && (
                        <div className="mt-4 p-3 bg-slate-800 rounded-lg flex items-center justify-between">
                          <div className="flex items-center">
                            <FileText className="w-5 h-5 mr-3 text-blue-500" />
                            <div>
                              <div className="text-sm font-medium text-white">
                                {importFile.name}
                              </div>
                              <div className="text-xs text-slate-400">
                                {(importFile.size / 1024 / 1024).toFixed(2)} MB
                              </div>
                            </div>
                          </div>
                          <button
                            onClick={() => setImportFile(null)}
                            className="p-1 text-slate-400 hover:text-red-400 transition-colors"
                          >
                            <X className="w-4 h-4" />
                          </button>
                        </div>
                      )}

                      <div className="mt-6 flex gap-3">
                        <button
                          onClick={() => setImportStep(1)}
                          className="flex-1 px-4 py-2 bg-slate-800 hover:bg-slate-700 text-white rounded-lg transition-colors"
                        >
                          Retour
                        </button>
                        <button
                          onClick={() => setImportStep(3)}
                          disabled={!importFile}
                          className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-slate-700 disabled:text-slate-500 text-white rounded-lg transition-colors"
                        >
                          Continuer
                        </button>
                      </div>
                    </div>
                  )}

                  {/* Step 3: Review and Import */}
                  {importStep === 3 && selectedSource && importFile && (
                    <div>
                      <h2 className="text-xl font-semibold text-white mb-4">
                        Vérifiez et importez
                      </h2>

                      <div className="space-y-4 mb-6">
                        <div className="p-4 bg-slate-800 rounded-lg">
                          <h3 className="text-sm font-medium text-white mb-3">
                            Détails de l'importation
                          </h3>
                          <div className="space-y-2 text-sm">
                            <div className="flex justify-between">
                              <span className="text-slate-400">Source:</span>
                              <span className="text-white">
                                {getSelectedSource()?.name}
                              </span>
                            </div>
                            <div className="flex justify-between">
                              <span className="text-slate-400">Fichier:</span>
                              <span className="text-white">
                                {importFile.name}
                              </span>
                            </div>
                            <div className="flex justify-between">
                              <span className="text-slate-400">Taille:</span>
                              <span className="text-white">
                                {(importFile.size / 1024 / 1024).toFixed(2)} MB
                              </span>
                            </div>
                          </div>
                        </div>

                        <div className="p-4 bg-blue-900 bg-opacity-20 border border-blue-800 rounded-lg">
                          <div className="flex items-start">
                            <Info className="w-5 h-5 mr-3 text-blue-500 mt-0.5 flex-shrink-0" />
                            <div className="text-sm text-slate-300">
                              <p className="font-medium text-white mb-1">
                                Important:
                              </p>
                              <ul className="list-disc list-inside space-y-1 text-slate-400">
                                <li>
                                  L'importation créera de nouveaux éléments dans
                                  votre coffre
                                </li>
                                <li>
                                  Les éléments existants avec les mêmes
                                  informations ne seront pas remplacés
                                </li>
                                <li>
                                  Les pièces jointes de fichiers ne sont pas
                                  importées
                                </li>
                                <li>
                                  Les Sends doivent être recréés manuellement
                                </li>
                              </ul>
                            </div>
                          </div>
                        </div>
                      </div>

                      <div className="flex gap-3">
                        <button
                          onClick={() => setImportStep(2)}
                          className="flex-1 px-4 py-2 bg-slate-800 hover:bg-slate-700 text-white rounded-lg transition-colors"
                        >
                          Retour
                        </button>
                        <button
                          onClick={handleImport}
                          disabled={isImporting}
                          className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-slate-700 disabled:text-slate-500 text-white rounded-lg transition-colors"
                        >
                          {isImporting ? (
                            <div className="flex items-center justify-center">
                              <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
                              Importation...
                            </div>
                          ) : (
                            "Importer les données"
                          )}
                        </button>
                      </div>
                    </div>
                  )}
                </div>
              )}

              {/* Import Results */}
              {importComplete && (
                <div className="space-y-6">
                  <div className="bg-slate-900 border border-slate-800 rounded-xl p-6">
                    <h2 className="text-xl font-semibold text-white mb-4">
                      Résultats de l'importation
                    </h2>

                    <div className="grid grid-cols-3 gap-4 mb-6">
                      <div className="text-center p-4 bg-slate-800 rounded-lg">
                        <div className="text-2xl font-bold text-green-500">
                          {importResults.successful}
                        </div>
                        <div className="text-xs text-slate-400">Réussis</div>
                      </div>
                      <div className="text-center p-4 bg-slate-800 rounded-lg">
                        <div className="text-2xl font-bold text-red-500">
                          {importResults.failed}
                        </div>
                        <div className="text-xs text-slate-400">Échoués</div>
                      </div>
                      <div className="text-center p-4 bg-slate-800 rounded-lg">
                        <div className="text-2xl font-bold text-yellow-500">
                          {importResults.skipped}
                        </div>
                        <div className="text-xs text-slate-400">Ignorés</div>
                      </div>
                    </div>

                    <div className="p-4 bg-green-900 bg-opacity-20 border border-green-800 rounded-lg">
                      <div className="flex items-center">
                        <CheckCircle className="w-5 h-5 mr-3 text-green-500" />
                        <div>
                          <p className="text-sm font-medium text-white">
                            Importation terminée!
                          </p>
                          <p className="text-xs text-slate-400">
                            Vos données ont été importées avec succès
                          </p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {/* Import Items Table */}
              {importItems.length > 0 && (
                <div className="bg-slate-900 border border-slate-800 rounded-xl">
                  <table className="w-full">
                    <thead className="bg-slate-800 sticky top-0">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                          Nom
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                          Type
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                          Source
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                          Statut
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                          Actions
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-800">
                      {filteredItems.map((item) => {
                        const ItemIcon = getItemIcon(item.type);
                        const StatusIcon = getStatusIcon(item.status);
                        return (
                          <tr
                            key={item.id}
                            className="hover:bg-slate-800 transition-colors"
                          >
                            <td className="px-6 py-4 whitespace-nowrap">
                              <div className="flex items-center">
                                <ItemIcon className="w-5 h-5 text-slate-400 mr-3" />
                                <div>
                                  <div className="text-sm font-medium text-white">
                                    {item.name}
                                  </div>
                                  <div className="text-xs text-slate-400 capitalize">
                                    {item.type}
                                  </div>
                                </div>
                              </div>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                              <div className="text-sm text-slate-300 capitalize">
                                {item.type === "login" && "Connexion"}
                                {item.type === "card" && "Carte"}
                                {item.type === "identity" && "Identité"}
                                {item.type === "secureNote" && "Note sécurisée"}
                              </div>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                              <div className="text-sm text-slate-300">
                                {importSources.find((s) => s.id === item.source)
                                  ?.name || item.source}
                              </div>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                              <div className="flex items-center">
                                <StatusIcon
                                  className={`w-4 h-4 mr-2 ${getStatusColor(item.status)}`}
                                />
                                <span
                                  className={`text-sm ${getStatusColor(item.status)}`}
                                >
                                  {item.status === "success" && "Réussi"}
                                  {item.status === "failed" && "Échoué"}
                                  {item.status === "skipped" && "Ignoré"}
                                  {item.status === "pending" && "En attente"}
                                </span>
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

              {/* Empty State */}
              {importItems.length === 0 && !importComplete && (
                <div className="flex flex-col items-center justify-center text-slate-400">
                  <Upload className="w-16 h-16 mb-4 text-slate-600" />
                  <h3 className="text-xl font-semibold mb-2">
                    Aucun élément à afficher
                  </h3>
                  <p className="text-sm mb-6">
                    Commencez par sélectionner une source et importer un fichier
                  </p>
                  <button
                    onClick={() => setImportStep(1)}
                    className="flex items-center px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                  >
                    <Plus className="w-5 h-5 mr-2" />
                    Commencer l'importation
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
