"use client";

import React, { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  Home,
  Key,
  Shield,
  Settings,
  FileText,
  Download,
  Upload,
  BarChart,
  Send,
  Clock,
  Wrench,
  User,
  Lock,
  CreditCard,
  Globe,
  ShieldCheck,
  HelpCircle,
} from "lucide-react";

const navigationItems = [
  { name: "Accueil", href: "/home", icon: Home },
  { name: "Coffre", href: "/vault", icon: Key },
  { name: "TOTP", href: "/totp", icon: Clock },
  { name: "Envoyer", href: "/sends", icon: Send },
  {
    name: "Outils",
    href: "/tools",
    icon: Wrench,
    subItems: [
      { name: "Générateur", href: "/tools/generator", icon: Shield },
      { name: "Importer", href: "/tools/import", icon: Download },
      { name: "Exporter", href: "/tools/export", icon: Upload },
    ],
  },
  { name: "Rapports", href: "/reports", icon: BarChart },
  {
    name: "Paramètres",
    href: "/settings",
    icon: Settings,
    subItems: [
      { name: "Compte", href: "/settings/account", icon: User },
      { name: "Sécurité", href: "/settings/security", icon: Lock },
      { name: "Préférences", href: "/settings/preferences", icon: Settings },
      { name: "Abonnement", href: "/settings/subscription", icon: CreditCard },
      {
        name: "Accès d'urgence",
        href: "/settings/emergency-access",
        icon: ShieldCheck,
      },
      {
        name: "Règles de domaine",
        href: "/settings/domain-rules",
        icon: Globe,
      },
    ],
  },
];

export function Sidebar() {
  const pathname = usePathname();
  const [expandedItems, setExpandedItems] = useState<string[]>([]);

  const toggleExpanded = (itemName: string) => {
    setExpandedItems((prev) =>
      prev.includes(itemName)
        ? prev.filter((item) => item !== itemName)
        : [...prev, itemName],
    );
  };

  return (
    <div className="w-64 bg-slate-900 border-r border-slate-800 flex flex-col">
      <div className="p-4 border-b border-slate-800">
        <div className="flex items-center space-x-3">
          <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
            <Shield className="w-5 h-5 text-white" />
          </div>
          <div>
            <h1 className="text-lg font-semibold text-white">Aether Vault</h1>
            <p className="text-xs text-slate-400">Gestionnaire de secrets</p>
          </div>
        </div>
      </div>

      <nav className="flex-1 p-4 space-y-1">
        {navigationItems.map((item) => {
          const isActive =
            pathname === item.href ||
            (item.href !== "/home" && pathname.startsWith(item.href));
          const isExpanded = expandedItems.includes(item.name);
          const hasSubItems = item.subItems && item.subItems.length > 0;

          return (
            <div key={item.name}>
              {hasSubItems ? (
                <button
                  onClick={() => toggleExpanded(item.name)}
                  className={`
                    w-full flex items-center px-3 py-2 text-sm font-medium rounded-lg transition-colors cursor-pointer
                    ${
                      isActive
                        ? "bg-blue-600 text-white"
                        : "text-slate-300 hover:bg-slate-800 hover:text-white"
                    }
                  `}
                >
                  <item.icon
                    className={`
                      mr-3 h-5 w-5 flex-shrink-0
                      ${isActive ? "text-white" : "text-slate-400"}
                    `}
                  />
                  <span className="flex-1 text-left">{item.name}</span>
                  <svg
                    className={`w-4 h-4 transition-transform ${isExpanded ? "rotate-90" : ""} ${
                      isActive ? "text-white" : "text-slate-400"
                    }`}
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 5l7 7-7 7"
                    />
                  </svg>
                </button>
              ) : (
                <Link
                  href={item.href}
                  className={`
                    flex items-center px-3 py-2 text-sm font-medium rounded-lg transition-colors
                    ${
                      isActive
                        ? "bg-blue-600 text-white"
                        : "text-slate-300 hover:bg-slate-800 hover:text-white"
                    }
                  `}
                >
                  <item.icon
                    className={`
                      mr-3 h-5 w-5 flex-shrink-0
                      ${isActive ? "text-white" : "text-slate-400"}
                    `}
                  />
                  <span className="flex-1 text-left">{item.name}</span>
                </Link>
              )}

              {hasSubItems && isExpanded && (
                <div className="mt-1 ml-4 space-y-1">
                  {item.subItems.map((subItem) => {
                    const isSubActive = pathname === subItem.href;

                    return (
                      <Link
                        key={subItem.name}
                        href={subItem.href}
                        className={`
                          flex items-center px-3 py-2 text-sm font-medium rounded-lg transition-colors
                          ${
                            isSubActive
                              ? "bg-blue-600 text-white"
                              : "text-slate-300 hover:bg-slate-800 hover:text-white"
                          }
                        `}
                      >
                        <subItem.icon
                          className={`
                            mr-3 h-4 w-4 flex-shrink-0
                            ${isSubActive ? "text-white" : "text-slate-400"}
                          `}
                        />
                        {subItem.name}
                      </Link>
                    );
                  })}
                </div>
              )}
            </div>
          );
        })}
      </nav>

      <div className="p-4 border-t border-slate-800">
        <div className="space-y-2">
          <button className="w-full flex items-center px-3 py-2 text-sm font-medium text-slate-300 rounded-lg hover:bg-slate-800 hover:text-white transition-colors">
            <HelpCircle className="mr-3 h-5 w-5 text-slate-400" />
            Aide & Support
          </button>
          <button className="w-full flex items-center px-3 py-2 text-sm font-medium text-slate-300 rounded-lg hover:bg-slate-800 hover:text-white transition-colors">
            <Download className="mr-3 h-5 w-5 text-slate-400" />
            Télécharger l'application
          </button>
        </div>
      </div>
    </div>
  );
}
