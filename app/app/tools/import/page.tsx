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
    }
  };

  const handleImport = async () => {
    if (!importFile || !selectedSource) return;

    setIsImporting(true);

    // Simulate import process
    await new Promise((resolve) => setTimeout(resolve, 2000));

    // Mock results
    setImportResults({
      successful: Math.floor(Math.random() * 50) + 10,
      failed: Math.floor(Math.random() * 5),
      skipped: Math.floor(Math.random() * 3),
    });

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
  };

  const getSelectedSource = () => {
    return importSources.find((source) => source.id === selectedSource);
  };

  return (
    <div className="h-full flex flex-col">
      <div className="flex-shrink-0 p-6 border-b border-slate-800">
        <div className="max-w-4xl mx-auto">
          <h1 className="text-3xl font-bold text-white mb-2">
            Importer des données
          </h1>
          <p className="text-slate-400">
            Importez vos mots de passe et données depuis d'autres gestionnaires
          </p>
        </div>
      </div>

      <div className="flex-1 overflow-auto">
        <div className="max-w-4xl mx-auto p-6">
          {/* Progress Steps */}
          <div className="mb-8">
            <div className="flex items-center justify-between">
              {importSteps.map((step, index) => (
                <div key={step.id} className="flex items-center flex-1">
                  <div className="flex items-center">
                    <div
                      className={`
                        w-10 h-10 rounded-full flex items-center justify-center text-sm font-medium transition-colors
                        ${
                          step.completed
                            ? "bg-blue-600 text-white"
                            : importStep >= step.id
                              ? "bg-slate-700 text-white"
                              : "bg-slate-800 text-slate-400"
                        }
                      `}
                    >
                      {step.completed ? (
                        <CheckCircle className="w-5 h-5" />
                      ) : (
                        step.id
                      )}
                    </div>
                    <div className="ml-3">
                      <h3
                        className={`text-sm font-medium ${
                          step.completed ? "text-white" : "text-slate-400"
                        }`}
                      >
                        {step.title}
                      </h3>
                      <p className="text-xs text-slate-500">
                        {step.description}
                      </p>
                    </div>
                  </div>
                  {index < importSteps.length - 1 && (
                    <div className="flex-1 mx-4">
                      <div
                        className={`h-px ${
                          step.completed ? "bg-blue-600" : "bg-slate-700"
                        }`}
                      />
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Main Content */}
            <div className="lg:col-span-2 space-y-6">
              {/* Step 1: Select Source */}
              {importStep === 1 && (
                <div className="bg-slate-900 border border-slate-800 rounded-lg p-6">
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
                          <Icon className={`w-6 h-6 mr-3 ${source.color}`} />
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
                <div className="bg-slate-900 border border-slate-800 rounded-lg p-6">
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
                        Cliquez pour sélectionner un fichier ou glissez-déposez
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
                <div className="bg-slate-900 border border-slate-800 rounded-lg p-6">
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
                          <span className="text-white">{importFile.name}</span>
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
                              Les éléments existants avec les mêmes informations
                              ne seront pas remplacés
                            </li>
                            <li>
                              Les pièces jointes de fichiers ne sont pas
                              importées
                            </li>
                            <li>Les Sends doivent être recréés manuellement</li>
                          </ul>
                        </div>
                      </div>
                    </div>
                  </div>

                  {!importComplete ? (
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
                  ) : (
                    <div className="space-y-4">
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

                      <div className="grid grid-cols-3 gap-4">
                        <div className="text-center p-3 bg-slate-800 rounded-lg">
                          <div className="text-2xl font-bold text-green-500">
                            {importResults.successful}
                          </div>
                          <div className="text-xs text-slate-400">Réussis</div>
                        </div>
                        <div className="text-center p-3 bg-slate-800 rounded-lg">
                          <div className="text-2xl font-bold text-red-500">
                            {importResults.failed}
                          </div>
                          <div className="text-xs text-slate-400">Échoués</div>
                        </div>
                        <div className="text-center p-3 bg-slate-800 rounded-lg">
                          <div className="text-2xl font-bold text-yellow-500">
                            {importResults.skipped}
                          </div>
                          <div className="text-xs text-slate-400">Ignorés</div>
                        </div>
                      </div>

                      <button
                        onClick={resetImport}
                        className="w-full px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                      >
                        Nouvelle importation
                      </button>
                    </div>
                  )}
                </div>
              )}
            </div>

            {/* Sidebar */}
            <div className="space-y-6">
              {/* Help */}
              <div className="bg-slate-900 border border-slate-800 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-white mb-4 flex items-center">
                  <Info className="w-5 h-5 mr-2 text-blue-500" />
                  Aide
                </h3>
                <div className="space-y-3 text-sm text-slate-400">
                  <p>
                    Avant d'importer, exportez vos données depuis votre
                    gestionnaire de mots de passe actuel.
                  </p>
                  <p>
                    Assurez-vous que le fichier est dans un format supporté et
                    n'est pas protégé par un mot de passe.
                  </p>
                  <p>
                    L'importation peut prendre quelques minutes selon la taille
                    de votre fichier.
                  </p>
                </div>
              </div>

              {/* Supported Formats */}
              <div className="bg-slate-900 border border-slate-800 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-white mb-4 flex items-center">
                  <FileText className="w-5 h-5 mr-2 text-green-500" />
                  Formats supportés
                </h3>
                <div className="space-y-2 text-sm">
                  <div className="flex items-center justify-between">
                    <span className="text-slate-400">CSV</span>
                    <span className="text-green-500">✓</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-slate-400">JSON</span>
                    <span className="text-green-500">✓</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-slate-400">XML</span>
                    <span className="text-green-500">✓</span>
                  </div>
                  <div className="flex items-center justify-between">
                    <span className="text-slate-400">1pif</span>
                    <span className="text-green-500">✓</span>
                  </div>
                </div>
              </div>

              {/* Tips */}
              <div className="bg-slate-900 border border-slate-800 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-white mb-4 flex items-center">
                  <AlertTriangle className="w-5 h-5 mr-2 text-yellow-500" />
                  Conseils
                </h3>
                <div className="space-y-3 text-sm text-slate-400">
                  <div className="flex items-start">
                    <ArrowRight className="w-4 h-4 mr-2 mt-0.5 text-blue-500 flex-shrink-0" />
                    <span>Sauvegardez votre coffre avant d'importer</span>
                  </div>
                  <div className="flex items-start">
                    <ArrowRight className="w-4 h-4 mr-2 mt-0.5 text-blue-500 flex-shrink-0" />
                    <span>Vérifiez les données après l'importation</span>
                  </div>
                  <div className="flex items-start">
                    <ArrowRight className="w-4 h-4 mr-2 mt-0.5 text-blue-500 flex-shrink-0" />
                    <span>Supprimez l'ancien coffre après vérification</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
