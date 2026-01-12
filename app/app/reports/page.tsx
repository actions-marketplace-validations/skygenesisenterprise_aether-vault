"use client";

import React, { useState } from "react";
import {
  FileText,
  Download,
  Calendar,
  Clock,
  TrendingUp,
  TrendingDown,
  Users,
  Shield,
  Lock,
  Key,
  CreditCard,
  User,
  AlertTriangle,
  CheckCircle,
  Info,
  Search,
  Filter,
  Plus,
  RefreshCw,
  Eye,
  Copy,
  Edit,
  Trash2,
  BarChart3,
  Activity,
  Zap,
  X,
} from "lucide-react";

interface Report {
  id: string;
  name: string;
  type: "security" | "usage" | "audit" | "compliance";
  description: string;
  generatedAt: Date;
  status: "completed" | "generating" | "failed";
  size: string;
  format: string;
  author: string;
}

interface ReportTemplate {
  id: string;
  name: string;
  type: "security" | "usage" | "audit" | "compliance";
  description: string;
  icon: React.ElementType;
  color: string;
  features: string[];
}

export default function ReportsPage() {
  const [selectedType, setSelectedType] = useState<string>("all");
  const [searchQuery, setSearchQuery] = useState("");
  const [showFavorites, setShowFavorites] = useState(false);
  const [showGenerateModal, setShowGenerateModal] = useState(false);
  const [selectedTemplate, setSelectedTemplate] = useState<string>("");
  const [isGenerating, setIsGenerating] = useState(false);

  const reports: Report[] = [
    {
      id: "1",
      name: "Rapport de sécurité - Q4 2024",
      type: "security",
      description: "Analyse complète des menaces et vulnérabilités",
      generatedAt: new Date("2024-12-15"),
      status: "completed",
      size: "2.4 MB",
      format: "PDF",
      author: "Admin",
    },
    {
      id: "2",
      name: "Rapport d'utilisation - Mensuel",
      type: "usage",
      description: "Statistiques d'utilisation du coffre",
      generatedAt: new Date("2024-12-01"),
      status: "completed",
      size: "1.1 MB",
      format: "Excel",
      author: "Système",
    },
    {
      id: "3",
      name: "Audit d'accès - Hebdomadaire",
      type: "audit",
      description: "Journal des accès et modifications",
      generatedAt: new Date("2024-12-08"),
      status: "generating",
      size: "-",
      format: "CSV",
      author: "Système",
    },
    {
      id: "4",
      name: "Rapport de conformité - GDPR",
      type: "compliance",
      description: "Vérification de conformité RGPD",
      generatedAt: new Date("2024-11-30"),
      status: "completed",
      size: "3.7 MB",
      format: "PDF",
      author: "Admin",
    },
    {
      id: "5",
      name: "Analyse des mots de passe faibles",
      type: "security",
      description: "Détection des mots de passe vulnérables",
      generatedAt: new Date("2024-12-10"),
      status: "failed",
      size: "-",
      format: "PDF",
      author: "Système",
    },
  ];

  const reportTemplates: ReportTemplate[] = [
    {
      id: "security-scan",
      name: "Analyse de sécurité",
      type: "security",
      description: "Scan complet des vulnérabilités et menaces",
      icon: Shield,
      color: "text-red-500",
      features: ["Vulnérabilités", "Menaces", "Recommandations"],
    },
    {
      id: "usage-stats",
      name: "Statistiques d'utilisation",
      type: "usage",
      description: "Métriques d'utilisation et tendances",
      icon: BarChart3,
      color: "text-blue-500",
      features: ["Connexions", "Utilisateurs actifs", "Tendances"],
    },
    {
      id: "access-audit",
      name: "Audit d'accès",
      type: "audit",
      description: "Journal détaillé des accès et modifications",
      icon: Activity,
      color: "text-green-500",
      features: ["Connexions", "Modifications", "Échecs"],
    },
    {
      id: "compliance-check",
      name: "Vérification de conformité",
      type: "compliance",
      description: "Contrôle des normes et réglementations",
      icon: CheckCircle,
      color: "text-purple-500",
      features: ["GDPR", "SOC2", "ISO27001"],
    },
    {
      id: "password-health",
      name: "Santé des mots de passe",
      type: "security",
      description: "Analyse de la force des mots de passe",
      icon: Lock,
      color: "text-orange-500",
      features: ["Force", "Répétitions", "Expirations"],
    },
  ];

  const itemTypes = [
    { value: "all", label: "Tous les rapports", icon: FileText },
    { value: "security", label: "Sécurité", icon: Shield },
    { value: "usage", label: "Utilisation", icon: BarChart3 },
    { value: "audit", label: "Audit", icon: Activity },
    { value: "compliance", label: "Conformité", icon: CheckCircle },
  ];

  const getReportIcon = (type: string) => {
    switch (type) {
      case "security":
        return Shield;
      case "usage":
        return BarChart3;
      case "audit":
        return Activity;
      case "compliance":
        return CheckCircle;
      default:
        return FileText;
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case "completed":
        return CheckCircle;
      case "generating":
        return RefreshCw;
      case "failed":
        return AlertTriangle;
      default:
        return Clock;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "completed":
        return "text-green-500";
      case "generating":
        return "text-blue-500";
      case "failed":
        return "text-red-500";
      default:
        return "text-slate-400";
    }
  };

  const getTypeColor = (type: string) => {
    switch (type) {
      case "security":
        return "text-red-500";
      case "usage":
        return "text-blue-500";
      case "audit":
        return "text-green-500";
      case "compliance":
        return "text-purple-500";
      default:
        return "text-slate-400";
    }
  };

  const handleGenerateReport = async () => {
    if (!selectedTemplate) return;

    setIsGenerating(true);
    setShowGenerateModal(false);

    // Simulate report generation
    await new Promise((resolve) => setTimeout(resolve, 3000));

    setIsGenerating(false);
    setSelectedTemplate("");
  };

  const handleDownloadReport = (reportId: string) => {
    // Simulate download
    console.log(`Downloading report ${reportId}`);
  };

  const filteredReports = reports.filter((report) => {
    const matchesSearch =
      report.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      report.description.toLowerCase().includes(searchQuery.toLowerCase()) ||
      report.author.toLowerCase().includes(searchQuery.toLowerCase());

    const matchesType = selectedType === "all" || report.type === selectedType;
    const matchesFavorites = !showFavorites || report.status === "completed";

    return matchesSearch && matchesType && matchesFavorites;
  });

  const stats = {
    total: reports.length,
    completed: reports.filter((r) => r.status === "completed").length,
    generating: reports.filter((r) => r.status === "generating").length,
    failed: reports.filter((r) => r.status === "failed").length,
  };

  return (
    <div className="h-full flex flex-col">
      <div className="flex-shrink-0 p-6 border-b border-slate-800">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-white mb-2">Rapports</h1>
            <p className="text-slate-400">
              Générez et consultez les rapports de sécurité et d'utilisation
            </p>
          </div>
          <button
            onClick={() => setShowGenerateModal(true)}
            className="flex items-center px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors"
          >
            <Plus className="w-5 h-5 mr-2" />
            Générer un rapport
          </button>
        </div>
      </div>

      <div className="flex-1 flex overflow-hidden">
        <div className="w-64 flex-shrink-0 border-r border-slate-800 overflow-y-auto">
          <div className="p-4 space-y-6">
            <div className="bg-slate-900 rounded-lg p-4">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-sm font-semibold text-white">Types</h3>
                <button className="text-slate-400 hover:text-white transition-colors">
                  <RefreshCw className="w-4 h-4" />
                </button>
              </div>

              <div className="space-y-4">
                <div>
                  <label className="block text-xs font-medium text-slate-400 mb-2">
                    Catégorie de rapports
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
                                ? "bg-slate-700 text-white border border-slate-600"
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
                    <CheckCircle className="w-4 h-4 mr-2" />
                    Terminés uniquement
                  </button>
                </div>
              </div>
            </div>

            {/* Statistics */}
            <div className="bg-slate-900 rounded-lg p-4">
              <h3 className="text-sm font-semibold text-white mb-4">
                Statistiques
              </h3>
              <div className="space-y-3">
                <div className="flex justify-between text-sm">
                  <span className="text-slate-400">Total</span>
                  <span className="text-white">{stats.total}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-slate-400">Terminés</span>
                  <span className="text-green-500">{stats.completed}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-slate-400">En cours</span>
                  <span className="text-blue-500">{stats.generating}</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-slate-400">Échoués</span>
                  <span className="text-red-500">{stats.failed}</span>
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
                  placeholder="Rechercher dans les rapports..."
                  className="w-full pl-10 pr-4 py-2 bg-slate-800 border border-slate-700 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-slate-600 focus:border-transparent"
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
              {/* Summary Cards */}
              <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                <div className="bg-slate-900 border border-slate-800 rounded-xl p-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm text-slate-400">Total</p>
                      <p className="text-2xl font-bold text-white">
                        {stats.total}
                      </p>
                    </div>
                    <FileText className="w-8 h-8 text-slate-500" />
                  </div>
                </div>
                <div className="bg-slate-900 border border-slate-800 rounded-xl p-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm text-slate-400">Terminés</p>
                      <p className="text-2xl font-bold text-green-500">
                        {stats.completed}
                      </p>
                    </div>
                    <CheckCircle className="w-8 h-8 text-green-500" />
                  </div>
                </div>
                <div className="bg-slate-900 border border-slate-800 rounded-xl p-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm text-slate-400">En cours</p>
                      <p className="text-2xl font-bold text-blue-500">
                        {stats.generating}
                      </p>
                    </div>
                    <RefreshCw className="w-8 h-8 text-blue-500" />
                  </div>
                </div>
                <div className="bg-slate-900 border border-slate-800 rounded-xl p-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <p className="text-sm text-slate-400">Échoués</p>
                      <p className="text-2xl font-bold text-red-500">
                        {stats.failed}
                      </p>
                    </div>
                    <AlertTriangle className="w-8 h-8 text-red-500" />
                  </div>
                </div>
              </div>

              {/* Reports Table */}
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
                        Date
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                        Statut
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                        Taille
                      </th>
                      <th className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                        Actions
                      </th>
                    </tr>
                  </thead>
                  <tbody className="divide-y divide-slate-800">
                    {filteredReports.map((report) => {
                      const ReportIcon = getReportIcon(report.type);
                      const StatusIcon = getStatusIcon(report.status);
                      return (
                        <tr
                          key={report.id}
                          className="hover:bg-slate-800 transition-colors"
                        >
                          <td className="px-6 py-4 whitespace-nowrap">
                            <div className="flex items-center">
                              <ReportIcon
                                className={`w-5 h-5 mr-3 ${getTypeColor(report.type)}`}
                              />
                              <div>
                                <div className="text-sm font-medium text-white">
                                  {report.name}
                                </div>
                                <div className="text-xs text-slate-400">
                                  {report.description}
                                </div>
                              </div>
                            </div>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            <div className="text-sm text-slate-300 capitalize">
                              {report.type === "security" && "Sécurité"}
                              {report.type === "usage" && "Utilisation"}
                              {report.type === "audit" && "Audit"}
                              {report.type === "compliance" && "Conformité"}
                            </div>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            <div className="text-sm text-slate-300">
                              {report.generatedAt.toLocaleDateString()}
                            </div>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            <div className="flex items-center">
                              <StatusIcon
                                className={`w-4 h-4 mr-2 ${getStatusColor(report.status)} ${report.status === "generating" ? "animate-spin" : ""}`}
                              />
                              <span
                                className={`text-sm ${getStatusColor(report.status)}`}
                              >
                                {report.status === "completed" && "Terminé"}
                                {report.status === "generating" && "En cours"}
                                {report.status === "failed" && "Échoué"}
                              </span>
                            </div>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            <div className="text-sm text-slate-300">
                              {report.size}
                            </div>
                          </td>
                          <td className="px-6 py-4 whitespace-nowrap">
                            <div className="flex items-center space-x-2">
                              <button
                                onClick={() => handleDownloadReport(report.id)}
                                disabled={report.status !== "completed"}
                                className="p-1 text-slate-400 hover:text-white transition-colors disabled:text-slate-600 disabled:cursor-not-allowed"
                              >
                                <Download className="w-4 h-4" />
                              </button>
                              <button className="p-1 text-slate-400 hover:text-white transition-colors">
                                <Eye className="w-4 h-4" />
                              </button>
                              <button className="p-1 text-slate-400 hover:text-white transition-colors">
                                <Copy className="w-4 h-4" />
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

              {/* Empty State */}
              {filteredReports.length === 0 && (
                <div className="flex flex-col items-center justify-center text-slate-400">
                  <FileText className="w-16 h-16 mb-4 text-slate-600" />
                  <h3 className="text-xl font-semibold mb-2">
                    Aucun rapport trouvé
                  </h3>
                  <p className="text-sm mb-6">
                    Commencez par générer votre premier rapport
                  </p>
                  <button
                    onClick={() => setShowGenerateModal(true)}
                    className="flex items-center px-4 py-2 bg-slate-700 hover:bg-slate-600 text-white rounded-lg transition-colors"
                  >
                    <Plus className="w-5 h-5 mr-2" />
                    Générer un rapport
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Generate Report Modal */}
      {showGenerateModal && (
        <div className="fixed inset-0 bg-black bg-opacity-30 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="bg-slate-900 border border-slate-800 rounded-xl p-6 w-full max-w-2xl max-h-[80vh] overflow-hidden">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-semibold text-white">
                Générer un rapport
              </h2>
              <button
                onClick={() => setShowGenerateModal(false)}
                className="text-slate-400 hover:text-white transition-colors"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-4">
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                {reportTemplates.map((template) => {
                  const Icon = template.icon;
                  return (
                    <button
                      key={template.id}
                      onClick={() => setSelectedTemplate(template.id)}
                      className={`
                        p-4 rounded-lg border-2 transition-all text-left
                        ${
                          selectedTemplate === template.id
                            ? "border-slate-600 bg-slate-800"
                            : "border-slate-700 hover:border-slate-600 hover:bg-slate-800"
                        }
                      `}
                    >
                      <div className="flex items-center mb-3">
                        <Icon className={`w-6 h-6 mr-3 ${template.color}`} />
                        <div className="font-medium text-white">
                          {template.name}
                        </div>
                      </div>
                      <div className="text-xs text-slate-400 mb-2">
                        {template.description}
                      </div>
                      <div className="flex flex-wrap gap-1">
                        {template.features.map((feature, index) => (
                          <span
                            key={index}
                            className="px-2 py-1 bg-slate-700 text-xs text-slate-300 rounded"
                          >
                            {feature}
                          </span>
                        ))}
                      </div>
                    </button>
                  );
                })}
              </div>

              <div className="pt-4 border-t border-slate-700">
                <div className="flex gap-3">
                  <button
                    onClick={() => setShowGenerateModal(false)}
                    className="flex-1 px-4 py-2 bg-slate-800 hover:bg-slate-700 text-white rounded-lg transition-colors"
                  >
                    Annuler
                  </button>
                  <button
                    onClick={handleGenerateReport}
                    disabled={!selectedTemplate || isGenerating}
                    className="flex-1 px-4 py-2 bg-slate-700 hover:bg-slate-600 disabled:bg-slate-800 disabled:text-slate-500 text-white rounded-lg transition-colors"
                  >
                    {isGenerating ? (
                      <div className="flex items-center justify-center">
                        <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2" />
                        Génération...
                      </div>
                    ) : (
                      "Générer"
                    )}
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
