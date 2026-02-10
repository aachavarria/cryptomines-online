import { useState, useEffect } from 'react';
import { combatReportsAPI, type CombatReport } from '../services/api';

export const useCombatReports = () => {
  const [reports, setReports] = useState<CombatReport[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchReports = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await combatReportsAPI.list();
      setReports(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch combat reports');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchReports();
  }, []);

  return {
    reports,
    loading,
    error,
    refetch: fetchReports,
  };
};

export const useCombatReport = (id: string | null) => {
  const [report, setReport] = useState<CombatReport | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) {
      setReport(null);
      return;
    }

    const fetchReport = async () => {
      try {
        setLoading(true);
        setError(null);
        const data = await combatReportsAPI.get(id);
        setReport(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch combat report');
      } finally {
        setLoading(false);
      }
    };

    fetchReport();
  }, [id]);

  return {
    report,
    loading,
    error,
  };
};
