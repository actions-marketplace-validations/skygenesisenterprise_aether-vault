"use client";

import React from "react";
import {
  Plus,
  Search,
  Shield,
  Clock,
  Users,
  Download,
  Upload,
  BarChart,
} from "lucide-react";
import Link from "next/link";

export default function HomePage() {
  return (
    <div className="p-6">
      <div className="max-w-7xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-white mb-2">
            Bienvenue dans Aether Vault
          </h1>
          <p className="text-slate-400">
            Gérez vos mots de passe, secrets et codes TOTP en toute sécurité
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <div className="bg-slate-900 border border-slate-800 rounded-lg p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="w-12 h-12 bg-blue-600 rounded-lg flex items-center justify-center">
                <Shield className="w-6 h-6 text-white" />
              </div>
              <span className="text-2xl font-bold text-white">0</span>
            </div>
            <h3 className="text-lg font-semibold text-white mb-1">
              Total des éléments
            </h3>
            <p className="text-sm text-slate-400">Dans votre coffre</p>
          </div>

          <div className="bg-slate-900 border border-slate-800 rounded-lg p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="w-12 h-12 bg-green-600 rounded-lg flex items-center justify-center">
                <Clock className="w-6 h-6 text-white" />
              </div>
              <span className="text-2xl font-bold text-white">0</span>
            </div>
            <h3 className="text-lg font-semibold text-white mb-1">
              Codes TOTP
            </h3>
            <p className="text-sm text-slate-400">
              Authentification à deux facteurs
            </p>
          </div>

          <div className="bg-slate-900 border border-slate-800 rounded-lg p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="w-12 h-12 bg-purple-600 rounded-lg flex items-center justify-center">
                <Users className="w-6 h-6 text-white" />
              </div>
              <span className="text-2xl font-bold text-white">0</span>
            </div>
            <h3 className="text-lg font-semibold text-white mb-1">
              Organisations
            </h3>
            <p className="text-sm text-slate-400">Partage d'équipe</p>
          </div>

          <div className="bg-slate-900 border border-slate-800 rounded-lg p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="w-12 h-12 bg-orange-600 rounded-lg flex items-center justify-center">
                <BarChart className="w-6 h-6 text-white" />
              </div>
              <span className="text-2xl font-bold text-white">0</span>
            </div>
            <h3 className="text-lg font-semibold text-white mb-1">Rapports</h3>
            <p className="text-sm text-slate-400">Analyse de sécurité</p>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2">
            <div className="bg-slate-900 border border-slate-800 rounded-lg p-6">
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-xl font-semibold text-white">
                  Éléments récents
                </h2>
                <Link
                  href="/vault"
                  className="text-blue-400 hover:text-blue-300 text-sm font-medium"
                >
                  Voir tout
                </Link>
              </div>

              <div className="space-y-4">
                <div className="flex items-center justify-center py-12 text-slate-400">
                  <div className="text-center">
                    <Search className="w-12 h-12 mx-auto mb-4 text-slate-600" />
                    <p className="text-lg font-medium mb-2">
                      Aucun élément trouvé
                    </p>
                    <p className="text-sm">
                      Commencez par ajouter votre premier élément au coffre
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div className="space-y-6">
            <div className="bg-slate-900 border border-slate-800 rounded-lg p-6">
              <h2 className="text-xl font-semibold text-white mb-4">
                Actions rapides
              </h2>
              <div className="space-y-3">
                <Link
                  href="/vault?action=new"
                  className="flex items-center w-full px-4 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                >
                  <Plus className="w-5 h-5 mr-3" />
                  Ajouter un élément
                </Link>

                <Link
                  href="/tools/generator"
                  className="flex items-center w-full px-4 py-3 bg-slate-800 hover:bg-slate-700 text-white rounded-lg transition-colors"
                >
                  <Shield className="w-5 h-5 mr-3" />
                  Générer un mot de passe
                </Link>

                <Link
                  href="/tools/import"
                  className="flex items-center w-full px-4 py-3 bg-slate-800 hover:bg-slate-700 text-white rounded-lg transition-colors"
                >
                  <Download className="w-5 h-5 mr-3" />
                  Importer des éléments
                </Link>
              </div>
            </div>

            <div className="bg-slate-900 border border-slate-800 rounded-lg p-6">
              <h2 className="text-xl font-semibold text-white mb-4">
                Conseils de sécurité
              </h2>
              <div className="space-y-3">
                <div className="flex items-start">
                  <div className="w-2 h-2 bg-green-500 rounded-full mt-2 mr-3 flex-shrink-0"></div>
                  <div>
                    <p className="text-sm text-white font-medium">
                      Activer l'authentification à deux facteurs
                    </p>
                    <p className="text-xs text-slate-400 mt-1">
                      Ajoutez une couche de sécurité supplémentaire
                    </p>
                  </div>
                </div>

                <div className="flex items-start">
                  <div className="w-2 h-2 bg-yellow-500 rounded-full mt-2 mr-3 flex-shrink-0"></div>
                  <div>
                    <p className="text-sm text-white font-medium">
                      Vérifier la force des mots de passe
                    </p>
                    <p className="text-xs text-slate-400 mt-1">
                      Assurez-vous que vos mots de passe sont robustes
                    </p>
                  </div>
                </div>

                <div className="flex items-start">
                  <div className="w-2 h-2 bg-blue-500 rounded-full mt-2 mr-3 flex-shrink-0"></div>
                  <div>
                    <p className="text-sm text-white font-medium">
                      Configurer l'accès d'urgence
                    </p>
                    <p className="text-xs text-slate-400 mt-1">
                      Permettez à un contact de confiance d'accéder à votre
                      compte
                    </p>
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
