"use client";

import { AreaChart, KPI, TrendChip } from "@heroui-pro/react";

const sparklineUp = [
	{ value: 30 },
	{ value: 35 },
	{ value: 28 },
	{ value: 42 },
	{ value: 38 },
	{ value: 45 },
	{ value: 50 },
	{ value: 48 },
	{ value: 55 },
	{ value: 60 },
	{ value: 58 },
	{ value: 65 },
];

const sparklineDown = [
	{ value: 65 },
	{ value: 60 },
	{ value: 62 },
	{ value: 55 },
	{ value: 58 },
	{ value: 52 },
	{ value: 50 },
	{ value: 48 },
	{ value: 45 },
	{ value: 42 },
	{ value: 44 },
	{ value: 40 },
];

export function Chart() {
	return (
		<div className="grid w-full max-w-[900px] grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
			<KPI>
				<KPI.Header>
					<KPI.Title>Total Revenue</KPI.Title>
				</KPI.Header>
				<KPI.Content className="grid-cols-[1fr_1fr] items-end">
					<div className="flex flex-col gap-1">
						<KPI.Value
							className="text-3xl"
							currency="USD"
							maximumFractionDigits={0}
							style="currency"
							value={228451}
						/>
						<TrendChip trend="up" variant="tertiary">
							3.3%
							<TrendChip.Suffix>last 30d</TrendChip.Suffix>
						</TrendChip>
					</div>
					<div className="min-w-0">
						<AreaChart
							data={sparklineUp}
							height={70}
							margin={{ bottom: 0, left: 0, right: 0, top: 4 }}
						>
							<defs>
								<linearGradient id="kpi-revenue" x1="0" x2="0" y1="0" y2="1">
									<stop
										offset="0%"
										stopColor="var(--chart-3)"
										stopOpacity={0.2}
									/>
									<stop
										offset="100%"
										stopColor="var(--chart-3)"
										stopOpacity={0.02}
									/>
								</linearGradient>
							</defs>
							<AreaChart.Area
								dataKey="value"
								dot={false}
								fill="url(#kpi-revenue)"
								stroke="var(--chart-3)"
								strokeWidth={1.5}
								type="monotone"
							/>
						</AreaChart>
					</div>
				</KPI.Content>
			</KPI>

			<KPI>
				<KPI.Header>
					<KPI.Title>Bounce Rate</KPI.Title>
				</KPI.Header>
				<KPI.Content className="grid-cols-[1fr_1fr] items-end">
					<div className="flex flex-col gap-1">
						<KPI.Value
							className="text-3xl"
							maximumFractionDigits={1}
							style="percent"
							value={0.423}
						/>
						<TrendChip trend="down" variant="tertiary">
							5.9%
							<TrendChip.Suffix>vs last 7d</TrendChip.Suffix>
						</TrendChip>
					</div>
					<div className="min-w-0">
						<AreaChart
							data={sparklineDown}
							height={70}
							margin={{ bottom: 0, left: 0, right: 0, top: 4 }}
						>
							<defs>
								<linearGradient id="kpi-bounce" x1="0" x2="0" y1="0" y2="1">
									<stop
										offset="0%"
										stopColor="var(--color-danger)"
										stopOpacity={0.2}
									/>
									<stop
										offset="100%"
										stopColor="var(--color-danger)"
										stopOpacity={0.02}
									/>
								</linearGradient>
							</defs>
							<AreaChart.Area
								dataKey="value"
								dot={false}
								fill="url(#kpi-bounce)"
								stroke="var(--color-danger)"
								strokeWidth={1.5}
								type="monotone"
							/>
						</AreaChart>
					</div>
				</KPI.Content>
			</KPI>

			<KPI>
				<KPI.Header>
					<KPI.Title>Active Users</KPI.Title>
				</KPI.Header>
				<KPI.Content className="grid-cols-[1fr_1fr] items-end">
					<div className="flex flex-col gap-1">
						<KPI.Value
							className="text-3xl"
							maximumFractionDigits={0}
							notation="compact"
							value={97859}
						/>
						<TrendChip trend="up" variant="tertiary">
							10.9%
							<TrendChip.Suffix>this month</TrendChip.Suffix>
						</TrendChip>
					</div>
					<div className="min-w-0">
						<AreaChart
							data={sparklineUp}
							height={70}
							margin={{ bottom: 0, left: 0, right: 0, top: 4 }}
						>
							<defs>
								<linearGradient id="kpi-users" x1="0" x2="0" y1="0" y2="1">
									<stop
										offset="0%"
										stopColor="var(--color-success)"
										stopOpacity={0.2}
									/>
									<stop
										offset="100%"
										stopColor="var(--color-success)"
										stopOpacity={0.02}
									/>
								</linearGradient>
							</defs>
							<AreaChart.Area
								dataKey="value"
								dot={false}
								fill="url(#kpi-users)"
								stroke="var(--color-success)"
								strokeWidth={1.5}
								type="monotone"
							/>
						</AreaChart>
					</div>
				</KPI.Content>
			</KPI>
		</div>
	);
}
