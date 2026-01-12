"use client";

import React, { useState, useEffect } from "react";
import {
  RefreshCw,
  Copy,
  Check,
  Shield,
  Lock,
  AlertTriangle,
  Clock,
  History,
  Settings,
  Zap,
  Eye,
  EyeOff,
  Save,
  Trash2,
  X,
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
  Edit,
} from "lucide-react";

interface GeneratorOptions {
  type: "password" | "passphrase" | "username";
  length: number;
  uppercase: boolean;
  lowercase: boolean;
  numbers: boolean;
  symbols: boolean;
  numWords: number;
  wordSeparator: string;
  capitalize: boolean;
  includeNumber: boolean;
  avoidAmbiguous: boolean;
}

interface GeneratedItem {
  id: string;
  value: string;
  type: string;
  timestamp: Date;
  strength?: number;
}

interface PasswordStrength {
  score: number;
  text: string;
  color: string;
  timeToCrack: string;
  entropy: number;
}

export default function GeneratorPage() {
  const [generatedValue, setGeneratedValue] = useState("");
  const [copied, setCopied] = useState(false);
  const [showHistory, setShowHistory] = useState(false);
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [history, setHistory] = useState<GeneratedItem[]>([]);
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedType, setSelectedType] = useState<string>("all");
  const [showFavorites, setShowFavorites] = useState(false);
  const [showPasswords, setShowPasswords] = useState<Set<string>>(new Set());
  const [options, setOptions] = useState<GeneratorOptions>({
    type: "password",
    length: 16,
    uppercase: true,
    lowercase: true,
    numbers: true,
    symbols: true,
    numWords: 4,
    wordSeparator: "-",
    capitalize: false,
    includeNumber: false,
    avoidAmbiguous: false,
  });
  const [passwordStrength, setPasswordStrength] = useState<PasswordStrength>({
    score: 0,
    text: "Très faible",
    color: "text-red-500",
    timeToCrack: "Instantané",
    entropy: 0,
  });

  const wordList = [
    "apple",
    "banana",
    "coffee",
    "dragon",
    "elephant",
    "forest",
    "garden",
    "house",
    "island",
    "jungle",
    "kitchen",
    "mountain",
    "ocean",
    "planet",
    "river",
    "sunset",
    "tiger",
    "umbrella",
    "valley",
    "window",
    "yellow",
    "zebra",
    "butterfly",
    "castle",
    "diamond",
    "eagle",
    "flower",
    "galaxy",
    "horizon",
    "lighthouse",
    "meadow",
    "nature",
    "oasis",
    "paradise",
    "quantum",
    "rainbow",
    "starlight",
    "thunder",
    "universe",
    "volcano",
    "waterfall",
    "crystal",
    "dream",
    "emerald",
    "fountain",
    "glacier",
    "harmony",
    "infinity",
    "journey",
    "kingdom",
    "lunar",
    "mystic",
    "nebula",
    "oracle",
    "phoenix",
    "quasar",
    "radiant",
    "serene",
    "twilight",
    "utopia",
    "victory",
    "wisdom",
    "zenith",
    "cosmos",
    "aurora",
    "bliss",
    "cascade",
    "dawn",
    "echo",
    "flame",
    "grove",
    "haven",
    "iris",
    "jewel",
    "kaleidoscope",
    "lagoon",
    "prism",
    "quill",
    "serenity",
    "temple",
    "valor",
    "whisper",
    "xenon",
    "yonder",
    "zephyr",
  ];

  const generatePassword = () => {
    let charset = "";

    if (options.lowercase) {
      charset += options.avoidAmbiguous
        ? "abcdefghijkmnopqrstuvwxyz"
        : "abcdefghijklmnopqrstuvwxyz";
    }
    if (options.uppercase) {
      charset += options.avoidAmbiguous
        ? "ABCDEFGHJKLMNPQRSTUVWXYZ"
        : "ABCDEFGHIJKLMNOPQRSTUVWXYZ";
    }
    if (options.numbers) {
      charset += options.avoidAmbiguous ? "23456789" : "0123456789";
    }
    if (options.symbols) {
      charset += "!@#$%^&*()_+-=[]{}|;:,.<>?";
    }

    if (charset === "") {
      setGeneratedValue("Veuillez sélectionner au moins un type de caractère");
      return;
    }

    let password = "";
    for (let i = 0; i < options.length; i++) {
      password += charset.charAt(Math.floor(Math.random() * charset.length));
    }

    setGeneratedValue(password);
    addToHistory(password, "password");
  };

  const generatePassphrase = () => {
    const words = [];
    const availableWords = [...wordList];

    for (let i = 0; i < options.numWords; i++) {
      const randomIndex = Math.floor(Math.random() * availableWords.length);
      let word = availableWords[randomIndex];
      availableWords.splice(randomIndex, 1);

      if (options.capitalize) {
        word = word.charAt(0).toUpperCase() + word.slice(1);
      }
      words.push(word);
    }

    if (options.includeNumber) {
      words.push(Math.floor(Math.random() * 9999) + 1);
    }

    const passphrase = words.join(options.wordSeparator);
    setGeneratedValue(passphrase);
    addToHistory(passphrase, "passphrase");
  };

  const generateUsername = () => {
    const adjectives = [
      "cool",
      "smart",
      "brave",
      "quick",
      "silent",
      "dark",
      "bright",
      "wild",
      "calm",
      "bold",
    ];
    const nouns = [
      "wolf",
      "eagle",
      "tiger",
      "lion",
      "dragon",
      "phoenix",
      "shadow",
      "storm",
      "blade",
      "hunter",
    ];
    const numbers = Math.floor(Math.random() * 9999) + 1;

    const adjective = adjectives[Math.floor(Math.random() * adjectives.length)];
    const noun = nouns[Math.floor(Math.random() * nouns.length)];

    const username = `${adjective}${noun}${numbers}`;
    setGeneratedValue(username);
    addToHistory(username, "username");
  };

  const generate = () => {
    switch (options.type) {
      case "password":
        generatePassword();
        break;
      case "passphrase":
        generatePassphrase();
        break;
      case "username":
        generateUsername();
        break;
    }
  };

  const calculateStrength = (value: string): PasswordStrength => {
    if (!value || value.includes("Veuillez sélectionner")) {
      return {
        score: 0,
        text: "Très faible",
        color: "text-red-500",
        timeToCrack: "Instantané",
        entropy: 0,
      };
    }

    let score = 0;
    let charsetSize = 0;

    if (/[a-z]/.test(value)) charsetSize += 26;
    if (/[A-Z]/.test(value)) charsetSize += 26;
    if (/[0-9]/.test(value)) charsetSize += 10;
    if (/[^a-zA-Z0-9]/.test(value)) charsetSize += 32;

    const entropy = value.length * Math.log2(charsetSize);

    if (value.length >= 8) score += 1;
    if (value.length >= 12) score += 1;
    if (value.length >= 16) score += 1;
    if (value.length >= 20) score += 1;

    if (/[a-z]/.test(value)) score += 1;
    if (/[A-Z]/.test(value)) score += 1;
    if (/[0-9]/.test(value)) score += 1;
    if (/[^a-zA-Z0-9]/.test(value)) score += 1;

    let strength: PasswordStrength;

    if (score <= 2) {
      strength = {
        score: 25,
        text: "Très faible",
        color: "text-red-500",
        timeToCrack: "Instantané",
        entropy,
      };
    } else if (score <= 4) {
      strength = {
        score: 50,
        text: "Faible",
        color: "text-orange-500",
        timeToCrack: "Quelques minutes",
        entropy,
      };
    } else if (score <= 6) {
      strength = {
        score: 75,
        text: "Bon",
        color: "text-yellow-500",
        timeToCrack: "Plusieurs années",
        entropy,
      };
    } else {
      strength = {
        score: 100,
        text: "Excellent",
        color: "text-green-500",
        timeToCrack: "Plusieurs siècles",
        entropy,
      };
    }

    return strength;
  };

  const addToHistory = (value: string, type: string) => {
    const newItem: GeneratedItem = {
      id: Date.now().toString(),
      value,
      type,
      timestamp: new Date(),
      strength:
        type === "password" ? calculateStrength(value).score : undefined,
    };

    setHistory((prev) => [newItem, ...prev.slice(0, 19)]);
  };

  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(generatedValue);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error("Failed to copy:", err);
    }
  };

  const clearHistory = () => {
    setHistory([]);
  };

  const regenerate = () => {
    generate();
  };

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
      case "password":
        return Lock;
      case "passphrase":
        return Shield;
      case "username":
        return User;
      default:
        return Key;
    }
  };

  const itemTypes = [
    { value: "all", label: "Tous les éléments", icon: Key },
    { value: "password", label: "Mots de passe", icon: Lock },
    { value: "passphrase", label: "Phrases secrètes", icon: Shield },
    { value: "username", label: "Noms d'utilisateur", icon: User },
  ];

  const filteredHistory = history.filter((item) => {
    const matchesSearch =
      item.value.toLowerCase().includes(searchQuery.toLowerCase()) ||
      item.type.toLowerCase().includes(searchQuery.toLowerCase());

    const matchesType = selectedType === "all" || item.type === selectedType;
    const matchesFavorites =
      !showFavorites || (item.strength && item.strength >= 75);

    return matchesSearch && matchesType && matchesFavorites;
  });

  useEffect(() => {
    generate();
  }, []);

  useEffect(() => {
    setPasswordStrength(calculateStrength(generatedValue));
  }, [generatedValue]);

  return (
    <div className="h-full flex flex-col">
      <div className="flex-shrink-0 p-6 border-b border-slate-800">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-white mb-2">Générateur</h1>
            <p className="text-slate-400">
              Créez des mots de passe et secrets sécurisés
            </p>
          </div>
          <button
            onClick={regenerate}
            className="flex items-center px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
          >
            <Plus className="w-5 h-5 mr-2" />
            Générer
          </button>
        </div>
      </div>

      <div className="flex-1 flex overflow-hidden">
        <div className="w-64 flex-shrink-0 border-r border-slate-800 overflow-y-auto">
          <div className="p-4 space-y-6">
            <div className="bg-slate-900 rounded-lg p-4">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-sm font-semibold text-white">Options</h3>
                <button className="text-slate-400 hover:text-white transition-colors">
                  <Settings className="w-4 h-4" />
                </button>
              </div>

              <div className="space-y-4">
                <div>
                  <label className="block text-xs font-medium text-slate-400 mb-2">
                    Type de génération
                  </label>
                  <div className="space-y-1">
                    {itemTypes.map((type) => {
                      const Icon = type.icon;
                      return (
                        <button
                          key={type.value}
                          onClick={() => {
                            setSelectedType(type.value);
                            if (type.value !== "all") {
                              setOptions({
                                ...options,
                                type: type.value as any,
                              });
                              setTimeout(generate, 100);
                            }
                          }}
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

                {options.type === "password" && (
                  <div className="space-y-3">
                    <div>
                      <label className="block text-xs font-medium text-slate-400 mb-2">
                        Longueur: {options.length}
                      </label>
                      <input
                        type="range"
                        min="4"
                        max="128"
                        value={options.length}
                        onChange={(e) => {
                          setOptions({
                            ...options,
                            length: parseInt(e.target.value),
                          });
                          setTimeout(generate, 100);
                        }}
                        className="w-full h-2 bg-slate-700 rounded-lg appearance-none cursor-pointer"
                      />
                    </div>

                    <div className="space-y-2">
                      {[
                        { key: "uppercase", label: "Majuscules" },
                        { key: "lowercase", label: "Minuscules" },
                        { key: "numbers", label: "Nombres" },
                        { key: "symbols", label: "Symboles" },
                      ].map(({ key, label }) => (
                        <label
                          key={key}
                          className="flex items-center space-x-2 text-sm text-slate-300"
                        >
                          <input
                            type="checkbox"
                            checked={
                              options[key as keyof GeneratorOptions] as boolean
                            }
                            onChange={(e) => {
                              setOptions({
                                ...options,
                                [key]: e.target.checked,
                              });
                              setTimeout(generate, 100);
                            }}
                            className="w-3 h-3 text-blue-600 bg-slate-800 border-slate-600 rounded"
                          />
                          {label}
                        </label>
                      ))}
                    </div>
                  </div>
                )}

                {options.type === "passphrase" && (
                  <div className="space-y-3">
                    <div>
                      <label className="block text-xs font-medium text-slate-400 mb-2">
                        Mots: {options.numWords}
                      </label>
                      <input
                        type="range"
                        min="3"
                        max="10"
                        value={options.numWords}
                        onChange={(e) => {
                          setOptions({
                            ...options,
                            numWords: parseInt(e.target.value),
                          });
                          setTimeout(generate, 100);
                        }}
                        className="w-full h-2 bg-slate-700 rounded-lg appearance-none cursor-pointer"
                      />
                    </div>

                    <div>
                      <label className="block text-xs font-medium text-slate-400 mb-2">
                        Séparateur
                      </label>
                      <select
                        value={options.wordSeparator}
                        onChange={(e) => {
                          setOptions({
                            ...options,
                            wordSeparator: e.target.value,
                          });
                          setTimeout(generate, 100);
                        }}
                        className="w-full px-2 py-1 bg-slate-800 border border-slate-700 rounded text-white text-sm"
                      >
                        <option value="-">Tiret</option>
                        <option value="_">Underscore</option>
                        <option value=" ">Espace</option>
                        <option value=".">Point</option>
                      </select>
                    </div>
                  </div>
                )}

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
                    Fort uniquement
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
                  placeholder="Rechercher dans l'historique..."
                  className="w-full pl-10 pr-4 py-2 bg-slate-800 border border-slate-700 rounded-lg text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
              <button
                onClick={() => setShowHistory(!showHistory)}
                className="flex items-center px-4 py-2 bg-slate-800 hover:bg-slate-700 text-white rounded-lg transition-colors"
              >
                <History className="w-5 h-5 mr-2" />
                Historique ({history.length})
              </button>
            </div>
          </div>

          <div className="flex-1 overflow-auto">
            <div className="p-6 space-y-6">
              {/* Generated Output */}
              <div className="bg-slate-900 border border-slate-800 rounded-xl p-6">
                <div className="mb-6">
                  <label className="block text-lg font-medium text-slate-300 mb-4">
                    {options.type === "password" && "Mot de passe généré"}
                    {options.type === "passphrase" && "Phrase secrète générée"}
                    {options.type === "username" && "Nom d'utilisateur généré"}
                  </label>

                  <div className="flex items-center space-x-3">
                    <div className="flex-1 relative">
                      <input
                        type={showAdvanced ? "text" : "password"}
                        value={generatedValue}
                        readOnly
                        className="w-full px-4 py-3 bg-slate-800 border border-slate-700 rounded-lg text-white font-mono focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      />
                      <button
                        onClick={() => setShowAdvanced(!showAdvanced)}
                        className="absolute right-3 top-1/2 transform -translate-y-1/2 p-2 text-slate-400 hover:text-white transition-colors"
                      >
                        {showAdvanced ? (
                          <EyeOff className="w-4 h-4" />
                        ) : (
                          <Eye className="w-4 h-4" />
                        )}
                      </button>
                    </div>
                    <button
                      onClick={copyToClipboard}
                      className="flex items-center px-4 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                    >
                      {copied ? (
                        <Check className="w-5 h-5" />
                      ) : (
                        <Copy className="w-5 h-5" />
                      )}
                    </button>
                    <button
                      onClick={regenerate}
                      className="flex items-center px-4 py-3 bg-slate-800 hover:bg-slate-700 text-white rounded-lg transition-colors"
                    >
                      <RefreshCw className="w-5 h-5" />
                    </button>
                  </div>
                </div>

                {/* Password Strength Analysis */}
                {options.type === "password" && (
                  <div className="space-y-4">
                    <div className="bg-slate-800 rounded-lg p-4">
                      <h3 className="text-sm font-semibold text-white mb-3">
                        Analyse de force
                      </h3>

                      <div className="space-y-3">
                        <div className="flex items-center justify-between">
                          <span className="text-xs text-slate-400">
                            Force du mot de passe
                          </span>
                          <span
                            className={`text-xs font-medium ${passwordStrength.color}`}
                          >
                            {passwordStrength.text}
                          </span>
                        </div>

                        <div className="h-2 bg-slate-700 rounded-full overflow-hidden">
                          <div
                            className={`h-full transition-all duration-500 ${
                              passwordStrength.score <= 25
                                ? "bg-red-500"
                                : passwordStrength.score <= 50
                                  ? "bg-orange-500"
                                  : passwordStrength.score <= 75
                                    ? "bg-yellow-500"
                                    : "bg-green-500"
                            }`}
                            style={{ width: `${passwordStrength.score}%` }}
                          />
                        </div>

                        <div className="grid grid-cols-2 gap-3 text-xs">
                          <div>
                            <span className="text-slate-400">
                              Temps pour craquer:
                            </span>
                            <span
                              className={`ml-1 font-medium ${passwordStrength.color}`}
                            >
                              {passwordStrength.timeToCrack}
                            </span>
                          </div>
                          <div>
                            <span className="text-slate-400">Entropie:</span>
                            <span
                              className={`ml-1 font-medium ${passwordStrength.color}`}
                            >
                              {passwordStrength.entropy.toFixed(1)} bits
                            </span>
                          </div>
                        </div>
                      </div>
                    </div>

                    {/* Character Analysis */}
                    <div className="bg-slate-800 rounded-lg p-4">
                      <h3 className="text-sm font-semibold text-white mb-3">
                        Composition
                      </h3>
                      <div className="grid grid-cols-2 gap-2 text-xs">
                        <div className="flex items-center justify-between p-2 bg-slate-700 rounded">
                          <span className="text-slate-300">Longueur</span>
                          <span className="text-white font-medium">
                            {generatedValue.length}
                          </span>
                        </div>
                        <div className="flex items-center justify-between p-2 bg-slate-700 rounded">
                          <span className="text-slate-300">Majuscules</span>
                          <span className="text-white font-medium">
                            {(generatedValue.match(/[A-Z]/g) || []).length}
                          </span>
                        </div>
                        <div className="flex items-center justify-between p-2 bg-slate-700 rounded">
                          <span className="text-slate-300">Minuscules</span>
                          <span className="text-white font-medium">
                            {(generatedValue.match(/[a-z]/g) || []).length}
                          </span>
                        </div>
                        <div className="flex items-center justify-between p-2 bg-slate-700 rounded">
                          <span className="text-slate-300">Nombres</span>
                          <span className="text-white font-medium">
                            {(generatedValue.match(/[0-9]/g) || []).length}
                          </span>
                        </div>
                        <div className="flex items-center justify-between p-2 bg-slate-700 rounded col-span-2">
                          <span className="text-slate-300">Symboles</span>
                          <span className="text-white font-medium">
                            {
                              (generatedValue.match(/[^a-zA-Z0-9]/g) || [])
                                .length
                            }
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                )}

                {/* Passphrase Info */}
                {options.type === "passphrase" && (
                  <div className="bg-slate-800 rounded-lg p-4">
                    <h3 className="text-sm font-semibold text-white mb-3">
                      Informations
                    </h3>
                    <div className="grid grid-cols-2 gap-3 text-xs">
                      <div className="flex items-center justify-between p-2 bg-slate-700 rounded">
                        <span className="text-slate-300">Mots</span>
                        <span className="text-white font-medium">
                          {options.numWords}
                        </span>
                      </div>
                      <div className="flex items-center justify-between p-2 bg-slate-700 rounded">
                        <span className="text-slate-300">Caractères</span>
                        <span className="text-white font-medium">
                          {generatedValue.length}
                        </span>
                      </div>
                      <div className="flex items-center justify-between p-2 bg-slate-700 rounded">
                        <span className="text-slate-300">Séparateur</span>
                        <span className="text-white font-medium">
                          "{options.wordSeparator}"
                        </span>
                      </div>
                      <div className="flex items-center justify-between p-2 bg-slate-700 rounded">
                        <span className="text-slate-300">Majuscules</span>
                        <span className="text-white font-medium">
                          {options.capitalize ? "Oui" : "Non"}
                        </span>
                      </div>
                    </div>
                  </div>
                )}

                {/* Username Info */}
                {options.type === "username" && (
                  <div className="bg-slate-800 rounded-lg p-4">
                    <h3 className="text-sm font-semibold text-white mb-3">
                      Informations
                    </h3>
                    <div className="grid grid-cols-2 gap-3 text-xs">
                      <div className="flex items-center justify-between p-2 bg-slate-700 rounded">
                        <span className="text-slate-300">Caractères</span>
                        <span className="text-white font-medium">
                          {generatedValue.length}
                        </span>
                      </div>
                      <div className="flex items-center justify-between p-2 bg-slate-700 rounded">
                        <span className="text-slate-300">Type</span>
                        <span className="text-white font-medium">
                          Aléatoire
                        </span>
                      </div>
                    </div>
                  </div>
                )}
              </div>

              {/* History Table */}
              {filteredHistory.length > 0 && (
                <div className="bg-slate-900 border border-slate-800 rounded-xl">
                  <table className="w-full">
                    <thead className="bg-slate-800 sticky top-0">
                      <tr>
                        <th className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                          Valeur
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                          Type
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                          Date
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                          Force
                        </th>
                        <th className="px-6 py-3 text-left text-xs font-medium text-slate-400 uppercase tracking-wider">
                          Actions
                        </th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-800">
                      {filteredHistory.map((item) => {
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
                                  <div className="text-sm font-medium text-white">
                                    {showPasswords.has(item.id)
                                      ? item.value
                                      : "•".repeat(item.value.length)}
                                  </div>
                                </div>
                              </div>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                              <div className="text-sm text-slate-300 capitalize">
                                {item.type === "password" && "Mot de passe"}
                                {item.type === "passphrase" && "Phrase secrète"}
                                {item.type === "username" &&
                                  "Nom d'utilisateur"}
                              </div>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                              <div className="text-sm text-slate-300">
                                {item.timestamp.toLocaleString()}
                              </div>
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                              {item.strength ? (
                                <div className="flex items-center">
                                  <div className="h-2 w-16 bg-slate-700 rounded-full overflow-hidden mr-2">
                                    <div
                                      className={`h-full ${
                                        item.strength >= 75
                                          ? "bg-green-500"
                                          : item.strength >= 50
                                            ? "bg-yellow-500"
                                            : "bg-red-500"
                                      }`}
                                      style={{ width: `${item.strength}%` }}
                                    />
                                  </div>
                                  <span
                                    className={`text-xs ${
                                      item.strength >= 75
                                        ? "text-green-500"
                                        : item.strength >= 50
                                          ? "text-yellow-500"
                                          : "text-red-500"
                                    }`}
                                  >
                                    {item.strength}%
                                  </span>
                                </div>
                              ) : (
                                <span className="text-sm text-slate-400">
                                  -
                                </span>
                              )}
                            </td>
                            <td className="px-6 py-4 whitespace-nowrap">
                              <div className="flex items-center space-x-2">
                                <button
                                  onClick={() =>
                                    togglePasswordVisibility(item.id)
                                  }
                                  className="p-1 text-slate-400 hover:text-white transition-colors"
                                >
                                  {showPasswords.has(item.id) ? (
                                    <EyeOff className="w-4 h-4" />
                                  ) : (
                                    <Eye className="w-4 h-4" />
                                  )}
                                </button>
                                <button
                                  onClick={async () => {
                                    await navigator.clipboard.writeText(
                                      item.value,
                                    );
                                  }}
                                  className="p-1 text-slate-400 hover:text-white transition-colors"
                                >
                                  <Copy className="w-4 h-4" />
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

              {filteredHistory.length === 0 && (
                <div className="flex flex-col items-center justify-center text-slate-400">
                  <Lock className="w-16 h-16 mb-4 text-slate-600" />
                  <h3 className="text-xl font-semibold mb-2">
                    Aucun élément trouvé
                  </h3>
                  <p className="text-sm mb-6">
                    Commencez par générer votre premier mot de passe
                  </p>
                  <button
                    onClick={regenerate}
                    className="flex items-center px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors"
                  >
                    <Plus className="w-5 h-5 mr-2" />
                    Générer
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
