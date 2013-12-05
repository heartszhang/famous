namespace famousfront.datamodels
{
  internal class FeedFlags
  {
    public const uint FeedFlagNone = 0;
    public const uint FeedFlagReaded = 1;
    public const uint FeedFlagStar = 1 << 1;
    public const uint FeedFlagSave = 1 << 2;
  }
}