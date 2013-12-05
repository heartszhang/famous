namespace famousfront.datamodels
{
  internal class FeedTypes
  {
    public const uint FeedTypeUnknown   = 0;
    public const uint FeedTypeRss       = 1 << 0;
    public const uint FeedTypeAtom      = 1 << 1;
    public const uint FeedTypeSinaweibo = 1 << 2;
    public const uint FeedTypeQqweibo   = 1 << 3;
    public const uint FeedTypeBlog      = 1 << 4;
    public const uint FeedTypeTweet     = 1 << 5;
    public const uint FeedTypeFeed      = FeedTypeRss | FeedTypeAtom;
  }
}