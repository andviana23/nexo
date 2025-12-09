import { MeiosPagamentoList } from '@/components/meios-pagamento/meios-pagamento-list';
import { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Tipos de Recebimento | NEXO',
  description: 'Gerencie os meios de pagamento e tipos de recebimento da sua barbearia',
};

export default function TiposRecebimentoPage() {
  return (
    <div className="space-y-6">
      <MeiosPagamentoList />
    </div>
  );
}
