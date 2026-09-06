from data_agent import run_data_agent
from quality_agent import run_quality_agent
from testing_agent import run_testing_agent


def checkpoint(name):
    print(f"\n========== CHECKPOINT: {name} ==========\n")


def main():
    print("Starting Multi-Agent Data Pipeline")

    # 1. Data Agent
    print("Starting Data Agent...")
    run_data_agent()

    checkpoint("DATA_AGENT_COMPLETE")

    # 2. Quality Agent
    print("Starting Quality Agent...")
    run_quality_agent()

    checkpoint("QUALITY_AGENT_COMPLETE")

    # 3. Testing Agent
    print("Starting Testing Agent...")
    run_testing_agent()

    checkpoint("TESTING_AGENT_COMPLETE")

    print("\nPipeline completed successfully!")


if __name__ == "__main__":
    main()