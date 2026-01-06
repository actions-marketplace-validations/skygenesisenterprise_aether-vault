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
    "meadow",
    "nebula",
    "oasis",
    "prism",
    "quill",
    "rainbow",
    "serenity",
    "temple",
    "utopia",
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

    // Calculate charset size
    if (/[a-z]/.test(value)) charsetSize += 26;
    if (/[A-Z]/.test(value)) charsetSize += 26;
    if (/[0-9]/.test(value)) charsetSize += 10;
    if (/[^a-zA-Z0-9]/.test(value)) charsetSize += 32;

    // Calculate entropy
    const entropy = value.length * Math.log2(charsetSize);

    // Length bonus
    if (value.length >= 8) score += 1;
    if (value.length >= 12) score += 1;
    if (value.length >= 16) score += 1;
    if (value.length >= 20) score += 1;

    // Character variety
    if (/[a-z]/.test(value)) score += 1;
    if (/[A-Z]/.test(value)) score += 1;
    if (/[0-9]/.test(value)) score += 1;
    if (/[^a-zA-Z0-9]/.test(value)) score += 1;

    // Determine strength
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

    setHistory((prev) => [newItem, ...prev.slice(0, 19)]); // Keep last 20 items
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

  useEffect(() => {
    generate();
  }, []);

  useEffect(() => {
    setPasswordStrength(calculateStrength(generatedValue));
  }, [generatedValue]);

  return (
    <div className="h-full flex flex-col bg-slate-950">
      {/* Header */}
      <div className="flex-shrink-0 p-6 border-b border-slate-800">
        <div className="max-w-7xl mx-auto">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-white mb-2">
                Générateur de sécurité
              </h1>
              <p className="text-slate-400">
                Créez des mots de passe, phrases secrètes et noms d'utilisateur
                sécurisés
              </p>
            </div>
            <div className="flex items-center space-x-3">
              <button
                onClick={() => setShowHistory(!showHistory)}
                className={`flex items-center px-4 py-2 rounded-lg transition-colors ${
                  showHistory
                    ? "bg-blue-600 text-white"
                    : "bg-slate-800 text-slate-300 hover:bg-slate-700"
                }`}
              >
                <History className="w-4 h-4 mr-2" />
                Historique ({history.length})
              </button>
              <button
                onClick={() => setShowAdvanced(!showAdvanced)}
                className={`flex items-center px-4 py-2 rounded-lg transition-colors ${
                  showAdvanced
                    ? "bg-purple-600 text-white"
                    : "bg-slate-800 text-slate-300 hover:bg-slate-700"
                }`}
              >
                <Settings className="w-4 h-4 mr-2" />
                Avancé
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 overflow-hidden">
        <div className="max-w-7xl mx-auto p-6 h-full">
          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 h-full">
            {/* Generated Output */}
            <div className="lg:col-span-2 space-y-6">
              {/* Main Generator Card */}
              <div className="bg-slate-900 border border-slate-800 rounded-xl p-8 h-full flex flex-col">
                <div className="mb-8">
                  <div className="flex items-center justify-between mb-4">
                    <label className="text-lg font-medium text-slate-300">
                      {options.type === "password" && "Mot de passe généré"}
                      {options.type === "passphrase" &&
                        "Phrase secrète générée"}
                      {options.type === "username" &&
                        "Nom d'utilisateur généré"}
                    </label>
                    <div className="flex items-center space-x-2">
                      <button
                        onClick={() => setShowAdvanced(!showAdvanced)}
                        className="p-2 text-slate-400 hover:text-white transition-colors"
                      >
                        <Settings className="w-4 h-4" />
                      </button>
                    </div>
                  </div>

                  <div className="flex items-center space-x-3">
                    <div className="flex-1 relative">
                      <input
                        type={showAdvanced ? "text" : "password"}
                        value={generatedValue}
                        readOnly
                        className="w-full px-6 py-4 bg-slate-800 border border-slate-700 rounded-xl text-white font-mono text-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
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
                      className="flex items-center px-6 py-4 bg-blue-600 hover:bg-blue-700 text-white rounded-xl transition-colors"
                    >
                      {copied ? (
                        <Check className="w-5 h-5" />
                      ) : (
                        <Copy className="w-5 h-5" />
                      )}
                    </button>
                    <button
                      onClick={regenerate}
                      className="flex items-center px-6 py-4 bg-slate-800 hover:bg-slate-700 text-white rounded-xl transition-colors"
                    >
                      <RefreshCw className="w-5 h-5" />
                    </button>
                  </div>
                </div>

                {/* Password Strength Analysis */}
                {options.type === "password" && (
                  <div className="space-y-6 flex-1">
                    <div className="bg-slate-800 rounded-xl p-6">
                      <h3 className="text-lg font-semibold text-white mb-4">
                        Analyse de force
                      </h3>

                      <div className="space-y-4">
                        <div className="flex items-center justify-between">
                          <span className="text-sm text-slate-400">
                            Force du mot de passe
                          </span>
                          <span
                            className={`text-sm font-medium ${passwordStrength.color}`}
                          >
                            {passwordStrength.text}
                          </span>
                        </div>

                        <div className="h-3 bg-slate-700 rounded-full overflow-hidden">
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

                        <div className="grid grid-cols-2 gap-4 text-sm">
                          <div>
                            <span className="text-slate-400">
                              Temps pour craquer:
                            </span>
                            <span
                              className={`ml-2 font-medium ${passwordStrength.color}`}
                            >
                              {passwordStrength.timeToCrack}
                            </span>
                          </div>
                          <div>
                            <span className="text-slate-400">Entropie:</span>
                            <span
                              className={`ml-2 font-medium ${passwordStrength.color}`}
                            >
                              {passwordStrength.entropy.toFixed(1)} bits
                            </span>
                          </div>
                        </div>
                      </div>
                    </div>

                    {/* Character Analysis */}
                    <div className="bg-slate-800 rounded-xl p-6">
                      <h3 className="text-lg font-semibold text-white mb-4">
                        Composition
                      </h3>
                      <div className="grid grid-cols-2 gap-3 text-sm">
                        <div className="flex items-center justify-between p-3 bg-slate-700 rounded-lg">
                          <span className="text-slate-300">Longueur</span>
                          <span className="text-white font-medium">
                            {generatedValue.length}
                          </span>
                        </div>
                        <div className="flex items-center justify-between p-3 bg-slate-700 rounded-lg">
                          <span className="text-slate-300">Majuscules</span>
                          <span className="text-white font-medium">
                            {(generatedValue.match(/[A-Z]/g) || []).length}
                          </span>
                        </div>
                        <div className="flex items-center justify-between p-3 bg-slate-700 rounded-lg">
                          <span className="text-slate-300">Minuscules</span>
                          <span className="text-white font-medium">
                            {(generatedValue.match(/[a-z]/g) || []).length}
                          </span>
                        </div>
                        <div className="flex items-center justify-between p-3 bg-slate-700 rounded-lg">
                          <span className="text-slate-300">Nombres</span>
                          <span className="text-white font-medium">
                            {(generatedValue.match(/[0-9]/g) || []).length}
                          </span>
                        </div>
                        <div className="flex items-center justify-between p-3 bg-slate-700 rounded-lg col-span-2">
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
                  <div className="space-y-6 flex-1">
                    <div className="bg-slate-800 rounded-xl p-6">
                      <h3 className="text-lg font-semibold text-white mb-4">
                        Informations
                      </h3>
                      <div className="grid grid-cols-2 gap-4 text-sm">
                        <div className="flex items-center justify-between p-3 bg-slate-700 rounded-lg">
                          <span className="text-slate-300">Mots</span>
                          <span className="text-white font-medium">
                            {options.numWords}
                          </span>
                        </div>
                        <div className="flex items-center justify-between p-3 bg-slate-700 rounded-lg">
                          <span className="text-slate-300">Caractères</span>
                          <span className="text-white font-medium">
                            {generatedValue.length}
                          </span>
                        </div>
                        <div className="flex items-center justify-between p-3 bg-slate-700 rounded-lg">
                          <span className="text-slate-300">Séparateur</span>
                          <span className="text-white font-medium">
                            "{options.wordSeparator}"
                          </span>
                        </div>
                        <div className="flex items-center justify-between p-3 bg-slate-700 rounded-lg">
                          <span className="text-slate-300">Majuscules</span>
                          <span className="text-white font-medium">
                            {options.capitalize ? "Oui" : "Non"}
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                )}

                {/* Username Info */}
                {options.type === "username" && (
                  <div className="space-y-6 flex-1">
                    <div className="bg-slate-800 rounded-xl p-6">
                      <h3 className="text-lg font-semibold text-white mb-4">
                        Informations
                      </h3>
                      <div className="grid grid-cols-2 gap-4 text-sm">
                        <div className="flex items-center justify-between p-3 bg-slate-700 rounded-lg">
                          <span className="text-slate-300">Caractères</span>
                          <span className="text-white font-medium">
                            {generatedValue.length}
                          </span>
                        </div>
                        <div className="flex items-center justify-between p-3 bg-slate-700 rounded-lg">
                          <span className="text-slate-300">Type</span>
                          <span className="text-white font-medium">
                            Aléatoire
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            </div>

            {/* Sidebar */}
            <div className="space-y-6 h-full">
              {/* Type Selection */}
              <div className="bg-slate-900 border border-slate-800 rounded-xl p-6">
                <h3 className="text-lg font-semibold text-white mb-4">
                  Type de génération
                </h3>
                <div className="space-y-3">
                  <button
                    onClick={() => {
                      setOptions({ ...options, type: "password" });
                      setTimeout(generate, 100);
                    }}
                    className={`w-full flex items-center p-4 rounded-xl border-2 transition-all ${
                      options.type === "password"
                        ? "border-blue-600 bg-blue-600 bg-opacity-20 text-white"
                        : "border-slate-700 text-slate-300 hover:border-slate-600 hover:text-white"
                    }`}
                  >
                    <Lock className="w-5 h-5 mr-3" />
                    <div className="text-left">
                      <div className="font-medium">Mot de passe</div>
                      <div className="text-xs opacity-75">
                        Caractères aléatoires
                      </div>
                    </div>
                  </button>

                  <button
                    onClick={() => {
                      setOptions({ ...options, type: "passphrase" });
                      setTimeout(generate, 100);
                    }}
                    className={`w-full flex items-center p-4 rounded-xl border-2 transition-all ${
                      options.type === "passphrase"
                        ? "border-blue-600 bg-blue-600 bg-opacity-20 text-white"
                        : "border-slate-700 text-slate-300 hover:border-slate-600 hover:text-white"
                    }`}
                  >
                    <Shield className="w-5 h-5 mr-3" />
                    <div className="text-left">
                      <div className="font-medium">Phrase secrète</div>
                      <div className="text-xs opacity-75">
                        Mots mémorisables
                      </div>
                    </div>
                  </button>

                  <button
                    onClick={() => {
                      setOptions({ ...options, type: "username" });
                      setTimeout(generate, 100);
                    }}
                    className={`w-full flex items-center p-4 rounded-xl border-2 transition-all ${
                      options.type === "username"
                        ? "border-blue-600 bg-blue-600 bg-opacity-20 text-white"
                        : "border-slate-700 text-slate-300 hover:border-slate-600 hover:text-white"
                    }`}
                  >
                    <History className="w-5 h-5 mr-3" />
                    <div className="text-left">
                      <div className="font-medium">Nom d'utilisateur</div>
                      <div className="text-xs opacity-75">
                        Identifiant unique
                      </div>
                    </div>
                  </button>
                </div>
              </div>

              {/* Options */}
              <div className="bg-slate-900 border border-slate-800 rounded-xl p-6 flex-1 overflow-y-auto">
                <h3 className="text-lg font-semibold text-white mb-4">
                  Options
                </h3>

                {options.type === "password" && (
                  <div className="space-y-6">
                    <div>
                      <label className="block text-sm font-medium text-slate-300 mb-3">
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
                      <div className="flex justify-between text-xs text-slate-400 mt-2">
                        <span>4</span>
                        <span>128</span>
                      </div>
                    </div>

                    <div className="space-y-3">
                      <label className="block text-sm font-medium text-slate-300 mb-3">
                        Caractères inclus
                      </label>

                      {[
                        {
                          key: "uppercase",
                          label: "Majuscules (A-Z)",
                          icon: "A",
                        },
                        {
                          key: "lowercase",
                          label: "Minuscules (a-z)",
                          icon: "a",
                        },
                        { key: "numbers", label: "Nombres (0-9)", icon: "9" },
                        {
                          key: "symbols",
                          label: "Symboles (!@#$%^&*)",
                          icon: "#",
                        },
                      ].map(({ key, label, icon }) => (
                        <label
                          key={key}
                          className="flex items-center space-x-3 p-3 bg-slate-800 rounded-lg cursor-pointer hover:bg-slate-700 transition-colors"
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
                            className="w-4 h-4 text-blue-600 bg-slate-800 border-slate-600 rounded focus:ring-blue-500"
                          />
                          <div className="w-8 h-8 bg-slate-700 rounded flex items-center justify-center text-sm font-mono text-slate-300">
                            {icon}
                          </div>
                          <span className="text-sm text-slate-300">
                            {label}
                          </span>
                        </label>
                      ))}

                      <label className="flex items-center space-x-3 p-3 bg-slate-800 rounded-lg cursor-pointer hover:bg-slate-700 transition-colors">
                        <input
                          type="checkbox"
                          checked={options.avoidAmbiguous}
                          onChange={(e) => {
                            setOptions({
                              ...options,
                              avoidAmbiguous: e.target.checked,
                            });
                            setTimeout(generate, 100);
                          }}
                          className="w-4 h-4 text-blue-600 bg-slate-800 border-slate-600 rounded focus:ring-blue-500"
                        />
                        <AlertTriangle className="w-5 h-5 text-yellow-500" />
                        <span className="text-sm text-slate-300">
                          Éviter les caractères ambigus (0O, 1l)
                        </span>
                      </label>
                    </div>
                  </div>
                )}

                {options.type === "passphrase" && (
                  <div className="space-y-6">
                    <div>
                      <label className="block text-sm font-medium text-slate-300 mb-3">
                        Nombre de mots: {options.numWords}
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
                      <div className="flex justify-between text-xs text-slate-400 mt-2">
                        <span>3</span>
                        <span>10</span>
                      </div>
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-slate-300 mb-3">
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
                        className="w-full px-3 py-2 bg-slate-800 border border-slate-700 rounded-lg text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      >
                        <option value="-">Tiret (-)</option>
                        <option value="_">Underscore (_)</option>
                        <option value=" ">Espace</option>
                        <option value=".">Point (.)</option>
                        <option value=",">Virgule (,)</option>
                      </select>
                    </div>

                    <div className="space-y-3">
                      <label className="flex items-center space-x-3 p-3 bg-slate-800 rounded-lg cursor-pointer hover:bg-slate-700 transition-colors">
                        <input
                          type="checkbox"
                          checked={options.capitalize}
                          onChange={(e) => {
                            setOptions({
                              ...options,
                              capitalize: e.target.checked,
                            });
                            setTimeout(generate, 100);
                          }}
                          className="w-4 h-4 text-blue-600 bg-slate-800 border-slate-600 rounded focus:ring-blue-500"
                        />
                        <span className="text-sm text-slate-300">
                          Majuscule au début de chaque mot
                        </span>
                      </label>

                      <label className="flex items-center space-x-3 p-3 bg-slate-800 rounded-lg cursor-pointer hover:bg-slate-700 transition-colors">
                        <input
                          type="checkbox"
                          checked={options.includeNumber}
                          onChange={(e) => {
                            setOptions({
                              ...options,
                              includeNumber: e.target.checked,
                            });
                            setTimeout(generate, 100);
                          }}
                          className="w-4 h-4 text-blue-600 bg-slate-800 border-slate-600 rounded focus:ring-blue-500"
                        />
                        <span className="text-sm text-slate-300">
                          Inclure un nombre à la fin
                        </span>
                      </label>
                    </div>
                  </div>
                )}

                {options.type === "username" && (
                  <div className="space-y-6">
                    <div className="p-4 bg-slate-800 rounded-lg">
                      <p className="text-sm text-slate-300">
                        Les noms d'utilisateur sont générés aléatoirement en
                        combinant des adjectifs, des noms et des nombres pour
                        garantir l'unicité.
                      </p>
                    </div>
                  </div>
                )}
              </div>

              {/* Security Tips */}
              <div className="bg-slate-900 border border-slate-800 rounded-xl p-6">
                <h3 className="text-lg font-semibold text-white mb-4 flex items-center">
                  <Zap className="w-5 h-5 mr-2 text-yellow-500" />
                  Conseils de sécurité
                </h3>
                <div className="space-y-3 text-sm text-slate-400">
                  <div className="flex items-start">
                    <Check className="w-4 h-4 mr-2 text-green-500 mt-0.5 flex-shrink-0" />
                    <span>Utilisez un mot de passe unique par compte</span>
                  </div>
                  <div className="flex items-start">
                    <Check className="w-4 h-4 mr-2 text-green-500 mt-0.5 flex-shrink-0" />
                    <span>Privilégiez les mots de passe de 16+ caractères</span>
                  </div>
                  <div className="flex items-start">
                    <Check className="w-4 h-4 mr-2 text-green-500 mt-0.5 flex-shrink-0" />
                    <span>Activez l'authentification à deux facteurs</span>
                  </div>
                  <div className="flex items-start">
                    <AlertTriangle className="w-4 h-4 mr-2 text-yellow-500 mt-0.5 flex-shrink-0" />
                    <span>Évitez les informations personnelles</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* History Modal */}
      {showHistory && (
        <div className="fixed inset-0 bg-black bg-opacity-50 backdrop-blur-sm flex items-center justify-center z-50">
          <div className="bg-slate-900 border border-slate-800 rounded-xl p-6 w-full max-w-2xl max-h-[80vh] overflow-hidden">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-xl font-semibold text-white">
                Historique des générations
              </h2>
              <div className="flex items-center space-x-2">
                <button
                  onClick={clearHistory}
                  className="flex items-center px-3 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors text-sm"
                >
                  <Trash2 className="w-4 h-4 mr-2" />
                  Vider
                </button>
                <button
                  onClick={() => setShowHistory(false)}
                  className="p-2 text-slate-400 hover:text-white transition-colors"
                >
                  <X className="w-5 h-5" />
                </button>
              </div>
            </div>

            <div className="overflow-y-auto max-h-[60vh]">
              {history.length === 0 ? (
                <div className="text-center py-12">
                  <Clock className="w-12 h-12 mx-auto mb-4 text-slate-600" />
                  <p className="text-slate-400">Aucun historique disponible</p>
                </div>
              ) : (
                <div className="space-y-2">
                  {history.map((item) => (
                    <div
                      key={item.id}
                      className="flex items-center justify-between p-4 bg-slate-800 rounded-lg hover:bg-slate-700 transition-colors"
                    >
                      <div className="flex-1">
                        <div className="flex items-center space-x-3">
                          <span className="text-xs px-2 py-1 bg-slate-700 rounded text-slate-300">
                            {item.type === "password" && "Mot de passe"}
                            {item.type === "passphrase" && "Phrase"}
                            {item.type === "username" && "Username"}
                          </span>
                          <code className="text-sm text-white font-mono">
                            {item.value}
                          </code>
                        </div>
                        <div className="text-xs text-slate-400 mt-1">
                          {item.timestamp.toLocaleString()}
                          {item.strength && (
                            <span
                              className={`ml-2 ${
                                item.strength >= 75
                                  ? "text-green-500"
                                  : item.strength >= 50
                                    ? "text-yellow-500"
                                    : "text-red-500"
                              }`}
                            >
                              Force: {item.strength}%
                            </span>
                          )}
                        </div>
                      </div>
                      <button
                        onClick={async () => {
                          await navigator.clipboard.writeText(item.value);
                        }}
                        className="p-2 text-slate-400 hover:text-white transition-colors"
                      >
                        <Copy className="w-4 h-4" />
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
